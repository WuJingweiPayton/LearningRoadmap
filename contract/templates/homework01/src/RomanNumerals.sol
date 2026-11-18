// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

/// @title Roman numeral conversions (LeetCode 12 / 13)
library RomanNumerals {
    function intToRoman(uint256 num) internal pure returns (string memory) {
        require(num > 0 && num <= 3999, "out of range");

        uint256[13] memory values = [
            uint256(1000),
            uint256(900),
            uint256(500),
            uint256(400),
            uint256(100),
            uint256(90),
            uint256(50),
            uint256(40),
            uint256(10),
            uint256(9),
            uint256(5),
            uint256(4),
            uint256(1)
        ];
        string[13] memory symbols =
            ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];

        bytes memory result;
        for (uint256 i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                num -= values[i];
                result = bytes.concat(result, bytes(symbols[i]));
            }
        }
        return string(result);
    }

    function romanToInt(string memory s) internal pure returns (uint256) {
        bytes memory chars = bytes(s);
        uint256 total;
        for (uint256 i = 0; i < chars.length; i++) {
            uint256 current = _romanValue(chars[i]);
            if (i + 1 < chars.length) {
                uint256 next = _romanValue(chars[i + 1]);
                if (current < next) {
                    total += next - current;
                    i++;
                    continue;
                }
            }
            total += current;
        }
        return total;
    }

    function _romanValue(bytes1 char) private pure returns (uint256) {
        if (char == "I") return 1;
        if (char == "V") return 5;
        if (char == "X") return 10;
        if (char == "L") return 50;
        if (char == "C") return 100;
        if (char == "D") return 500;
        if (char == "M") return 1000;
        revert("invalid roman");
    }
}

contract RomanNumeralsFacade {
    function intToRoman(uint256 num) external pure returns (string memory) {
        return RomanNumerals.intToRoman(num);
    }

    function romanToInt(string calldata s) external pure returns (uint256) {
        return RomanNumerals.romanToInt(s);
    }
}
