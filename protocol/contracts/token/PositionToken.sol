// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract PositionToken is ERC20, Ownable {
    mapping(address => bool) public isMinter;

    constructor(
        string memory name_,
        string memory symbol_
    ) ERC20(name_, symbol_) Ownable(msg.sender) {}

    function setMinter(address account, bool allowed) external onlyOwner {
        isMinter[account] = allowed;
    }

    function mint(address to, uint256 amount) external returns (bool) {
        require(isMinter[msg.sender], "caller is not minter");
        _mint(to, amount);
        return true;
    }

    function burn(address from, uint256 amount) external returns (bool) {
        require(isMinter[msg.sender], "caller is not minter"); // only minter can burn
        _burn(from, amount);
        return true;
    }
}