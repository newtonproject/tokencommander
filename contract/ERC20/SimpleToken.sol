pragma solidity ^0.5.0;

import "./ERC20.sol";
import "./ERC20Detailed.sol";

/**
 * @title SimpleToken
 * @dev Very simple ERC20 Token example, where all tokens are pre-assigned to the creator.
 * Note they can later distribute these tokens as they wish using `transfer` and other
 * `ERC20` functions.
 */
contract SimpleToken is ERC20, ERC20Detailed {

    /**
     * @dev Constructor that gives msg.sender all of existing tokens.
     */
    constructor(string memory _name, string memory _symbol, uint8 _decimals, uint256 _initialsupply) public
        ERC20Detailed(_name, _symbol, _decimals) {
        _mint(msg.sender, _initialsupply);
      }
}
