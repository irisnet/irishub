// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract Token is ERC20, ERC20Permit, Ownable, ReentrancyGuard {
    event SwapToNative(address from, string to, uint256 amount);

    uint8 private _scale;

    constructor(
        string memory name_,
        string memory symbol_,
        uint8 scale_
    ) ERC20(name_, symbol_) ERC20Permit(name_) Ownable(msg.sender) {
        _scale = scale_;
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
    function swapToNative(string memory to, uint256 amount)
        public
        nonReentrant
    {
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
}
