#!/usr/bin/env python3

import json
import pathlib
import re
import subprocess
import tempfile
import unittest


SCRIPT = pathlib.Path(__file__).with_name("export-genesis.sh")


def sanitizer(marker):
    source = SCRIPT.read_text()
    match = re.search(r"<<'%s'\n(.*?)\n%s$" % (marker, marker), source,
                      re.MULTILINE | re.DOTALL)
    if not match:
        raise AssertionError("sanitizer %s not found" % marker)
    return match.group(1)


def run_sanitizer(marker, value, extra=None, check=True):
    with tempfile.TemporaryDirectory() as temp:
        temp = pathlib.Path(temp)
        source = temp / "source.json"
        output = temp / "output.json"
        source.write_text(json.dumps(value))
        args = ["python3", "-", str(source)]
        if extra is not None:
            extra_path = temp / "extra.json"
            extra_path.write_text(json.dumps(extra))
            args.append(str(extra_path))
        args.append(str(output))
        result = subprocess.run(args, input=sanitizer(marker), text=True,
                                capture_output=True, check=check)
        fixed = json.loads(output.read_text()) if output.exists() else None
        return result, fixed


def run_id_sanitizer(module, value, check=True):
    with tempfile.TemporaryDirectory() as temp:
        temp = pathlib.Path(temp)
        source = temp / "source.json"
        output = temp / "output.json"
        source.write_text(json.dumps(value))
        result = subprocess.run(
            ["python3", "-", "sanitize-ids", str(source), str(output), module],
            input=sanitizer("PY"), text=True, capture_output=True, check=check)
        fixed = json.loads(output.read_text()) if output.exists() else None
        return result, fixed


def shell_function(name):
    source = SCRIPT.read_text()
    match = re.search(r"^%s\(\) \{\n.*?^\}\n" % re.escape(name), source,
                      re.MULTILINE | re.DOTALL)
    if not match:
        raise AssertionError("shell function %s not found" % name)
    return match.group(0)


def run_status_shell(body):
    helpers = "\n".join(shell_function(name) for name in
                        ("export_failed", "die", "export_success"))
    script = ("set -euo pipefail\nEXPORT_PHASE=init\n" + helpers
              + "\ntrap export_failed EXIT\n" + body)
    return subprocess.run(["bash"], input=script, text=True,
                          capture_output=True, check=False)


class BankSanitizerTest(unittest.TestCase):
    def test_sorts_units_without_rebuilding_valid_metadata(self):
        metadata = [{
            "base": "utoken",
            "display": "token",
            "name": "Token",
            "symbol": "TOKEN",
            "denom_units": [
                {"denom": "token", "exponent": 6, "aliases": ["coin"]},
                {"denom": "utoken", "exponent": 0, "aliases": ["micro"]},
            ],
        }]
        result, fixed = run_sanitizer("PY2", metadata, {"tokens": []})
        self.assertEqual("1", result.stdout.strip())
        self.assertEqual(["utoken", "token"],
                         [u["denom"] for u in fixed[0]["denom_units"]])
        self.assertEqual(["micro"], fixed[0]["denom_units"][0]["aliases"])

    def test_keeps_base_and_display_when_exponents_collide(self):
        metadata = [{
            "base": "utoken",
            "display": "token",
            "name": "Token",
            "symbol": "TOKEN",
            "denom_units": [
                {"denom": "utoken", "exponent": 0, "aliases": []},
                {"denom": "other", "exponent": 6, "aliases": ["old"]},
                {"denom": "token", "exponent": 6, "aliases": ["coin"]},
            ],
        }]
        result, fixed = run_sanitizer("PY2", metadata, {"tokens": []})
        self.assertEqual("2", result.stdout.strip())
        self.assertEqual(["utoken", "token"],
                         [u["denom"] for u in fixed[0]["denom_units"]])
        self.assertEqual("token", fixed[0]["display"])
        self.assertEqual(["coin"], fixed[0]["denom_units"][1]["aliases"])

        again, fixed_again = run_sanitizer("PY2", fixed, {"tokens": []})
        self.assertEqual("0", again.stdout.strip())
        self.assertEqual(fixed, fixed_again)

    def test_keeps_base_when_exponent_zero_collides(self):
        metadata = [{
            "base": "utoken",
            "display": "token",
            "name": "Token",
            "symbol": "TOKEN",
            "denom_units": [
                {"denom": "other", "exponent": 0, "aliases": ["old"]},
                {"denom": "utoken", "exponent": 0, "aliases": ["micro"]},
                {"denom": "token", "exponent": 6, "aliases": []},
            ],
        }]
        result, fixed = run_sanitizer("PY2", metadata, {"tokens": []})
        self.assertEqual("2", result.stdout.strip())
        self.assertEqual(["utoken", "token"],
                         [u["denom"] for u in fixed[0]["denom_units"]])
        self.assertEqual(["micro"], fixed[0]["denom_units"][0]["aliases"])


class GovSanitizerTest(unittest.TestCase):
    def test_repairs_deposit_and_expedited_voting_period(self):
        params = {
            "min_deposit": [{"denom": "unyan", "amount": "10000000"}],
            "expedited_min_deposit": [{"denom": "uiris", "amount": "50000000"}],
            "voting_period": "13m20s",
            "expedited_voting_period": "24h0m0s",
        }
        result, fixed = run_sanitizer("PY3", params)
        self.assertEqual("2", result.stdout.strip())
        self.assertEqual([{"denom": "unyan", "amount": "50000000"}],
                         fixed["expedited_min_deposit"])
        self.assertEqual("400s", fixed["expedited_voting_period"])

    def test_retains_valid_params(self):
        params = {
            "min_deposit": [{"denom": "unyan", "amount": "10000000"}],
            "expedited_min_deposit": [{"denom": "unyan", "amount": "50000000"}],
            "voting_period": "13m20s",
            "expedited_voting_period": "5m0s",
        }
        result, fixed = run_sanitizer("PY3", params)
        self.assertEqual("0", result.stdout.strip())
        self.assertEqual(params, fixed)

    def test_halves_fractional_duration_at_nanosecond_precision(self):
        params = {
            "min_deposit": [],
            "voting_period": "1.000000003s",
            "expedited_voting_period": "2s",
        }
        result, fixed = run_sanitizer("PY3", params)
        self.assertEqual("1", result.stdout.strip())
        self.assertEqual("0.500000001s", fixed["expedited_voting_period"])

    def test_rejects_non_positive_regular_voting_period(self):
        params = {
            "min_deposit": [],
            "voting_period": "0s",
            "expedited_voting_period": "24h0m0s",
        }
        result, fixed = run_sanitizer("PY3", params, check=False)
        self.assertNotEqual(0, result.returncode)
        self.assertIsNone(fixed)
        self.assertIn("voting_period must be positive", result.stderr)

    def test_rejects_period_too_short_for_positive_expedited_period(self):
        params = {
            "min_deposit": [],
            "voting_period": "1ns",
            "expedited_voting_period": "1ns",
        }
        result, fixed = run_sanitizer("PY3", params, check=False)
        self.assertNotEqual(0, result.returncode)
        self.assertIsNone(fixed)
        self.assertIn("too short to derive a positive", result.stderr)


class AssetIdSanitizerTest(unittest.TestCase):
    def test_nft_drops_overlong_ids_and_keeps_boundary_and_tibc_denom(self):
        keep = "a" * 101
        long_id = "b" * 102
        tibc = "tibc-" + "c" * 120
        state = {"collections": [
            {"nfts": [{"id": "validToken", "owner": "owner"}],
             "denom": {"id": long_id, "name": "long"}, "extra": {"x": 1}},
            {"denom": {"id": keep}, "nfts": [{"id": keep, "owner": "owner"}]},
            {"denom": {"id": tibc}, "nfts": [{"id": "tibcToken", "owner": "owner"}]},
            {"denom": {"id": "tokens"}, "nfts": [{"id": long_id, "owner": "owner"}]},
        ]}
        result, fixed = run_id_sanitizer("nft", state)
        self.assertEqual(0, result.returncode)
        self.assertEqual([keep, tibc, "tokens"],
                         [c["denom"]["id"] for c in fixed["collections"]])
        self.assertEqual(keep, fixed["collections"][0]["nfts"][0]["id"])
        self.assertEqual("tibcToken", fixed["collections"][1]["nfts"][0]["id"])
        self.assertEqual([], fixed["collections"][2]["nfts"])

    def test_overlong_tokens_drop_without_cross_collection_collision(self):
        token = "z" * 102
        state = {"collections": [
            {"denom": {"id": "alpha"}, "nfts": [{"id": token, "owner": "one"}]},
            {"nfts": [{"id": token, "owner": "two"}], "denom": {"id": "beta"}},
        ]}
        _, fixed = run_id_sanitizer("nft", state)
        self.assertEqual([], fixed["collections"][0]["nfts"])
        self.assertEqual([], fixed["collections"][1]["nfts"])

    def test_mt_keeps_digit_leading_native_hex_ids_and_field_order(self):
        denom = "259edd57e552854d42bc4a0d98dc7a48fddeae343ad428c9df1d4b09e0ab525a"
        mt_id = "9" + "b" * 63
        boundary = "A" * 101
        state = {
            "owners": [{"address": "owner", "denoms": [{
                "balances": [{"amount": "7", "mt_id": mt_id, "memo": "keep"}],
                "denom_id": denom,
            }]}],
            "collections": [{
                "mts": [{"id": mt_id, "supply": "7", "data": "AA=="}],
                "denom": {"id": denom, "name": "D"},
            }, {"denom": {"id": boundary},
                "mts": [{"id": boundary, "supply": "0"}]}],
        }
        _, fixed = run_id_sanitizer("mt", state)
        owner_denom = fixed["owners"][0]["denoms"][0]
        self.assertEqual(denom, fixed["collections"][0]["denom"]["id"])
        self.assertEqual(mt_id, fixed["collections"][0]["mts"][0]["id"])
        self.assertEqual(denom, owner_denom["denom_id"])
        self.assertEqual(mt_id, owner_denom["balances"][0]["mt_id"])
        self.assertEqual("7", owner_denom["balances"][0]["amount"])
        self.assertEqual("7", fixed["collections"][0]["mts"][0]["supply"])
        self.assertEqual("keep", owner_denom["balances"][0]["memo"])
        self.assertEqual(boundary, fixed["collections"][1]["denom"]["id"])
        self.assertEqual(boundary, fixed["collections"][1]["mts"][0]["id"])

        _, again = run_id_sanitizer("mt", fixed)
        self.assertEqual(fixed, again)

    def test_invalid_overlong_denoms_drop_without_collision(self):
        prefix = "c" * 101
        state = {"collections": [
            {"denom": {"id": prefix + "x"}, "nfts": []},
            {"denom": {"id": prefix + "y"}, "nfts": []},
        ]}
        result, fixed = run_id_sanitizer("nft", state)
        self.assertEqual([], fixed["collections"])
        self.assertIn("2 invalid collection(s) dropped", result.stderr)

    def test_overlong_token_drops_without_collision_with_valid_token(self):
        prefix = "q" * 101
        state = {"collections": [{"denom": {"id": "denom"}, "nfts": [
            {"id": prefix + "x", "owner": "one"},
            {"id": prefix, "owner": "two"},
        ]}]}
        _, fixed = run_id_sanitizer("nft", state)
        self.assertEqual([prefix],
                         [nft["id"] for nft in fixed["collections"][0]["nfts"]])

    def test_empty_state_and_unknown_mt_references_are_not_invented(self):
        _, empty = run_id_sanitizer("mt", {"owners": [], "collections": []})
        self.assertEqual({"owners": [], "collections": []}, empty)

        unknown_denom = "unknownDenom"
        unknown_mt = "unknownToken"
        state = {"owners": [{"denoms": [{
            "denom_id": unknown_denom,
            "balances": [{"mt_id": unknown_mt, "amount": "9"}],
        }]}], "collections": []}
        _, fixed = run_id_sanitizer("mt", state)
        self.assertEqual(state, fixed)

        _, nulls = run_id_sanitizer("mt", {"owners": None,
                                                   "collections": None})
        self.assertEqual({"owners": None, "collections": None}, nulls)
        _, null_assets = run_id_sanitizer(
            "nft", {"collections": [{"denom": {"id": "valid"},
                                       "nfts": None}]})
        self.assertIsNone(null_assets["collections"][0]["nfts"])

    def test_nft_drops_invalid_character_ids_in_any_array_position(self):
        reported = "9«Dã¦ïOF­%`U-%"
        state = {"collections": [
            {"denom": {"id": "mixed"}, "nfts": [
                {"id": reported, "owner": "drop-first", "data": "bad"},
                {"id": "validToken", "owner": "keep", "data": "payload"},
                {"id": "bad%middle", "owner": "drop-middle"},
                {"id": "validLast", "owner": "keep-last"},
                {"id": "bad-last!", "owner": "drop-last"},
            ]},
            {"denom": {"id": "allbad"}, "nfts": [
                {"id": "%bad", "owner": "one"},
                {"id": "also_bad", "owner": "two"},
            ]},
        ]}
        result, fixed = run_id_sanitizer("nft", state)
        self.assertIn(json.dumps(reported), result.stderr)
        self.assertEqual(["validToken", "validLast"],
                         [nft["id"] for nft in fixed["collections"][0]["nfts"]])
        self.assertEqual("payload", fixed["collections"][0]["nfts"][0]["data"])
        self.assertEqual([], fixed["collections"][1]["nfts"])
        self.assertIn("0 invalid collection(s) dropped, 5 NFT record(s) dropped",
                      result.stderr)

        _, again = run_id_sanitizer("nft", fixed)
        self.assertEqual(fixed, again)

    def test_nft_filters_original_ids_without_collision_checks(self):
        invalid_prefix = "bad%" + "x" * 97
        valid_boundary = "s" * 101
        state = {"collections": [{"denom": {"id": "denom"}, "nfts": [
            {"id": invalid_prefix + "a", "owner": "drop-one"},
            {"id": invalid_prefix + "b", "owner": "drop-two"},
            {"id": valid_boundary + "tail", "owner": "drop-long"},
        ]}]}
        _, fixed = run_id_sanitizer("nft", state)
        self.assertEqual([], fixed["collections"][0]["nfts"])

        collision = {"collections": [{"denom": {"id": "denom"}, "nfts": [
            {"id": valid_boundary + "%", "owner": "drop-invalid"},
            {"id": valid_boundary, "owner": "keep-valid"},
        ]}]}
        _, fixed = run_id_sanitizer("nft", collision)
        self.assertEqual([valid_boundary],
                         [nft["id"] for nft in fixed["collections"][0]["nfts"]])

    def test_invalid_denom_drops_collection_and_all_nested_records(self):
        state = {"collections": [
            {"denom": {"id": "bad_denom", "name": "drop"}, "nfts": [
                {"id": "validOne", "owner": "one"},
                {"id": "bad%two", "owner": "two"},
            ]},
            {"denom": {"id": "valid"}, "nfts": []},
        ]}
        result, fixed = run_id_sanitizer("nft", state)
        self.assertEqual(["valid"],
                         [c["denom"]["id"] for c in fixed["collections"]])
        self.assertIn("collection denom", result.stderr)
        self.assertIn("1 invalid collection(s) dropped, 2 NFT record(s) dropped",
                      result.stderr)

    def test_mt_drops_invalid_definitions_and_matching_owner_refs(self):
        long_denom = "d" * 102
        long_mt = "m" * 102
        state = {
            "owners": [{"address": "owner", "denoms": [
                {"denom_id": "bad_denom", "balances": [
                    {"mt_id": "nestedToken", "amount": "1"}]},
                {"denom_id": long_denom, "balances": [
                    {"mt_id": "nestedToken", "amount": "2"}]},
                {"balances": [
                    {"mt_id": "bad%token", "amount": "3"},
                    {"mt_id": "validToken", "amount": "7", "memo": "keep"},
                    {"mt_id": long_mt, "amount": "4"},
                ], "denom_id": "validDenom"},
            ]}, {"address": "empty", "denoms": [{
                "denom_id": "bad_denom",
                "balances": [{"mt_id": "nestedToken", "amount": "5"}],
            }]}],
            "collections": [
                {"denom": {"id": "bad_denom"},
                 "mts": [{"id": "nestedToken", "supply": "1"}]},
                {"denom": {"id": long_denom},
                 "mts": [{"id": "nestedToken", "supply": "2"}]},
                {"mts": [
                    {"id": "bad%token", "supply": "3"},
                    {"id": "validToken", "supply": "7", "data": "keep"},
                    {"id": long_mt, "supply": "4"},
                 ], "denom": {"id": "validDenom"}},
            ],
        }
        result, fixed = run_id_sanitizer("mt", state)
        self.assertEqual(["validDenom"],
                         [c["denom"]["id"] for c in fixed["collections"]])
        self.assertEqual(["validToken"],
                         [mt["id"] for mt in fixed["collections"][0]["mts"]])
        self.assertEqual("7", fixed["collections"][0]["mts"][0]["supply"])
        self.assertEqual("keep", fixed["collections"][0]["mts"][0]["data"])
        denoms = fixed["owners"][0]["denoms"]
        self.assertEqual(["validDenom"], [d["denom_id"] for d in denoms])
        self.assertEqual(["validToken"],
                         [b["mt_id"] for b in denoms[0]["balances"]])
        self.assertEqual("7", denoms[0]["balances"][0]["amount"])
        self.assertEqual("keep", denoms[0]["balances"][0]["memo"])
        self.assertEqual({"address": "empty", "denoms": []}, fixed["owners"][1])
        self.assertIn("2 invalid collection(s) dropped, 4 MT record(s) dropped",
                      result.stderr)
        self.assertIn("3 owner denom reference(s) dropped", result.stderr)
        self.assertIn("5 owner balance reference(s) dropped", result.stderr)

    def test_same_valid_asset_id_in_different_collections_is_retained(self):
        state = {"collections": [
            {"denom": {"id": "alpha"},
             "nfts": [{"id": "sameToken", "owner": "one"}]},
            {"denom": {"id": "beta"},
             "nfts": [{"id": "sameToken", "owner": "two"}]},
        ]}
        _, fixed = run_id_sanitizer("nft", state)
        self.assertEqual(["sameToken", "sameToken"],
                         [c["nfts"][0]["id"] for c in fixed["collections"]])

    def test_asset_logs_escape_terminal_control_characters(self):
        invalid = ("\x00\x1f\x7f\x80\x85\x9b2J\x9d0;hidden\x9c"
                   "\u2028\u2029")
        metadata = "owner\x1b]0;hidden\x07\u2028"
        nft_state = {"collections": [{"denom": {"id": "valid"}, "nfts": [
            {"id": invalid, "owner": metadata, "data": metadata},
        ]}]}
        mt_state = {
            "collections": [{"denom": {"id": "valid"}, "mts": [
                {"id": invalid, "supply": "9", "data": metadata},
            ]}],
            "owners": [{"address": metadata, "denoms": [{
                "denom_id": "valid",
                "balances": [{"mt_id": invalid, "amount": "9"}],
            }]}],
        }
        for module, state in (("nft", nft_state), ("mt", mt_state)):
            with self.subTest(module=module):
                result, fixed = run_id_sanitizer(module, state)
                unsafe = [ch for ch in result.stderr
                          if ch != "\n" and not 0x20 <= ord(ch) <= 0x7e]
                self.assertEqual([], unsafe)
                self.assertIn(json.dumps(invalid), result.stderr)
                assets = fixed["collections"][0]["nfts" if module == "nft"
                                                         else "mts"]
                self.assertEqual([], assets)
                self.assertIn("1 NFT record(s) dropped" if module == "nft"
                              else "1 MT record(s) dropped", result.stderr)
                if module == "mt":
                    self.assertEqual([], fixed["owners"][0]["denoms"][0]
                                     ["balances"])
                    self.assertIn("1 owner balance reference(s) dropped",
                                  result.stderr)


class ExportStatusTest(unittest.TestCase):
    def test_failed_child_reports_phase_and_original_exit_code(self):
        result = run_status_shell("""
json_tool() { return 137; }
EXPORT_PHASE=sanitize_nft
json_tool sanitize-ids source output nft || die "nft ID sanitize failed"
export_success chain 1 denom 2 passed
""")
        self.assertEqual(137, result.returncode)
        self.assertNotIn("EXPORT_SUCCESS", result.stdout + result.stderr)
        self.assertIn("nft ID sanitize failed", result.stderr)
        self.assertIn("EXPORT_FAILED phase=sanitize_nft exit_code=137",
                      result.stderr)

    def test_success_marker_reports_validation_status(self):
        for validation in ("passed", "skipped"):
            with self.subTest(validation=validation):
                result = run_status_shell(
                    "export_success nyancat-10 30846001 unyan "
                    "1000390968884229927 %s\n" % validation)
                self.assertEqual(0, result.returncode)
                self.assertEqual("", result.stderr)
                self.assertEqual(
                    "EXPORT_SUCCESS chain_id=nyancat-10 "
                    "initial_height=30846001 denom=unyan "
                    "supply=1000390968884229927 module_validation=%s\n"
                    % validation,
                    result.stdout)


if __name__ == "__main__":
    unittest.main()
