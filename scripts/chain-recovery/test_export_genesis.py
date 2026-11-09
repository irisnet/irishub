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
    def test_nft_truncates_overlong_ids_and_keeps_boundary_and_tibc_denom(self):
        keep = "a" * 101
        long_id = "b" * 102
        tibc = "tibc-" + "c" * 120
        state = {"collections": [
            {"nfts": [{"id": long_id, "owner": "owner", "data": "payload"}],
             "denom": {"id": long_id, "name": "long"}, "extra": {"x": 1}},
            {"denom": {"id": keep}, "nfts": [{"id": keep, "owner": "owner"}]},
            {"denom": {"id": tibc}, "nfts": []},
        ]}
        result, fixed = run_id_sanitizer("nft", state)
        self.assertEqual(0, result.returncode)
        self.assertEqual(long_id[:101], fixed["collections"][0]["denom"]["id"])
        self.assertEqual(long_id[:101], fixed["collections"][0]["nfts"][0]["id"])
        self.assertEqual("payload", fixed["collections"][0]["nfts"][0]["data"])
        self.assertEqual({"x": 1}, fixed["collections"][0]["extra"])
        self.assertEqual(keep, fixed["collections"][1]["denom"]["id"])
        self.assertEqual(tibc, fixed["collections"][2]["denom"]["id"])

    def test_same_truncated_token_id_is_allowed_in_different_collections(self):
        token = "z" * 102
        state = {"collections": [
            {"denom": {"id": "alpha"}, "nfts": [{"id": token, "owner": "one"}]},
            {"nfts": [{"id": token, "owner": "two"}], "denom": {"id": "beta"}},
        ]}
        _, fixed = run_id_sanitizer("nft", state)
        self.assertEqual(token[:101], fixed["collections"][0]["nfts"][0]["id"])
        self.assertEqual(token[:101], fixed["collections"][1]["nfts"][0]["id"])

    def test_mt_updates_definitions_and_references_without_changing_amounts(self):
        denom = "d" * 102
        mt_id = "m" * 102
        state = {
            "owners": [{"address": "owner", "denoms": [{
                "balances": [{"amount": "7", "mt_id": mt_id, "memo": "keep"}],
                "denom_id": denom,
            }]}],
            "collections": [{
                "mts": [{"id": mt_id, "supply": "7", "data": "AA=="}],
                "denom": {"id": denom, "name": "D"},
            }],
        }
        _, fixed = run_id_sanitizer("mt", state)
        owner_denom = fixed["owners"][0]["denoms"][0]
        self.assertEqual(denom[:101], fixed["collections"][0]["denom"]["id"])
        self.assertEqual(mt_id[:101], fixed["collections"][0]["mts"][0]["id"])
        self.assertEqual(denom[:101], owner_denom["denom_id"])
        self.assertEqual(mt_id[:101], owner_denom["balances"][0]["mt_id"])
        self.assertEqual("7", owner_denom["balances"][0]["amount"])
        self.assertEqual("7", fixed["collections"][0]["mts"][0]["supply"])
        self.assertEqual("keep", owner_denom["balances"][0]["memo"])

        _, again = run_id_sanitizer("mt", fixed)
        self.assertEqual(fixed, again)

    def test_rejects_long_long_denom_collision_without_output(self):
        prefix = "c" * 101
        state = {"collections": [
            {"denom": {"id": prefix + "x"}, "nfts": []},
            {"denom": {"id": prefix + "y"}, "nfts": []},
        ]}
        result, fixed = run_id_sanitizer("nft", state, check=False)
        self.assertNotEqual(0, result.returncode)
        self.assertIsNone(fixed)
        self.assertIn("collision", result.stderr)

    def test_rejects_long_short_token_collision_without_output(self):
        prefix = "q" * 101
        state = {"collections": [{"denom": {"id": "denom"}, "nfts": [
            {"id": prefix + "x", "owner": "one"},
            {"id": prefix, "owner": "two"},
        ]}]}
        result, fixed = run_id_sanitizer("nft", state, check=False)
        self.assertNotEqual(0, result.returncode)
        self.assertIsNone(fixed)
        self.assertIn("collision", result.stderr)

    def test_empty_state_and_unknown_mt_references_are_not_invented(self):
        _, empty = run_id_sanitizer("mt", {"owners": [], "collections": []})
        self.assertEqual({"owners": [], "collections": []}, empty)

        unknown_denom = "u" * 102
        unknown_mt = "v" * 102
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


if __name__ == "__main__":
    unittest.main()
