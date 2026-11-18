// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @notice Simple counter for homework05 abigen / Go binding demo.
contract Counter {
    uint256 public count;

    function increment() external {
        count += 1;
    }
}