// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/introspection/ERC165.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";

contract Token is ERC20, Ownable, ReentrancyGuard, Initializable, ERC165 {
    event SwapToNative(address from, string to, uint256 amount);

    uint8 private _scale;
    uint8 public constant VERSION = 1;

    constructor(
        string memory name_,
        string memory symbol_,
        uint8 scale_
    ) ERC20(name_, symbol_) Ownable(msg.sender) initializer {
        _scale = scale_;
    }

    /**
     * @dev Sets the values for {name}, {symbol},{decimals} and {owner}.
     *
     * these values can only be set once during construction or initialize.
     */
    function initialize(
        string memory name_,
        string memory symbol_,
        uint8 scale_,
        address owner_
    ) public initializer {
        _name = name_;
        _symbol = symbol_;
        _scale = scale_;
        _transferOwnership(owner_);
    }

    /**
     * @dev Returns the number of decimals used to get its user representation.
     * For example, if `decimals` equals `2`, a balance of `505` tokens should
     * be displayed to a user as `5.05` (`505 / 10 ** 2`).
     *
     * Tokens usually opt for a value of 18, imitating the relationship between
     * Ether and Wei. This is the default value returned by this function, unless
     * it's overridden.
     *
     * NOTE: This information is only used for _display_ purposes: it in
     * no way affects any of the arithmetic of the contract, including
     * {IERC20-balanceOf} and {IERC20-transfer}.
     */
    function decimals() public view override returns (uint8) {
        return _scale;
    }

    /**
     * @dev Creates a `amount` amount of tokens and assigns them to `account`, by transferring it from address(0).
     * Relies on the `_update` mechanism
     *
     * Emits a {Transfer} event with `from` set to the zero address.
     *
     * NOTE: This function is not virtual, {_update} should be overridden instead.
     */
    function mint(address account, uint256 amount) public onlyOwner {
        _mint(account, amount);
    }

    /**
     * @dev Destroys a `amount` amount of tokens from `account`, lowering the total supply.
     * Relies on the `_update` mechanism.
     *
     * Emits a {Transfer} event with `to` set to the zero address.
     *
     * NOTE: This function is not virtual, {_update} should be overridden instead
     */
    function burn(address account, uint256 amount) public onlyOwner {
        _burn(account, amount);
    }

    /**
     *
     * Requirements:
     *
     * - `to` cannot be the zero address.
     * - `amount` caller must have a balance of at least `amount`.
     */
    function swapToNative(
        string memory to,
        uint256 amount
    ) public nonReentrant {
        require(bytes(to).length > 0, "to must be vaild iaa address");

        address sender = _msgSender();
        _burn(sender, amount);
        emit SwapToNative(sender, to, amount);
    }

    /**
     *
     * Requirements:
     *
     * - `from` authorizer address.
     * - `to` cannot be the zero address.
     * - `amount` from must have a balance of at least `amount`.
     */
    function swapToNativeFrom(
        address from,
        string memory to,
        uint256 amount
    ) public nonReentrant {
        require(bytes(to).length > 0, "to must be vaild iaa address");

        address spender = _msgSender();
        _spendAllowance(from, spender, amount);

        _burn(from, amount);
        emit SwapToNative(from, to, amount);
    }

    /**
     * @dev See {IERC165-supportsInterface}
     */
    function supportsInterface(
        bytes4 interfaceId
    ) public view override returns (bool) {
        return
            interfaceId == type(IERC20).interfaceId ||
            interfaceId == type(IERC20Metadata).interfaceId ||
            super.supportsInterface(interfaceId);
    }
}
