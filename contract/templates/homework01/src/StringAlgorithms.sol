// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

/// @title String algorithms (homework 01)
library StringAlgorithms {
    function reverse(string memory input) internal pure returns (string memory) {
        bytes memory chars = bytes(input);
        bytes memory reversed = new bytes(chars.length);
        for (uint256 i = 0; i < chars.length; i++) {
            reversed[i] = chars[chars.length - 1 - i];
        }
        return string(reversed);
    }
}

contract StringReverse {
    function reverse(string calldata input) external pure returns (string memory) {
        return StringAlgorithms.reverse(input);
    }
}
