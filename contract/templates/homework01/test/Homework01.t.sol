// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Voting} from "../src/Voting.sol";
import {StringReverse} from "../src/StringAlgorithms.sol";
import {RomanNumeralsFacade} from "../src/RomanNumerals.sol";
import {ArrayAlgorithmsFacade} from "../src/ArrayAlgorithms.sol";

contract Homework01Test is Test {
    Voting voting;
    StringReverse stringReverse;
    RomanNumeralsFacade roman;
    ArrayAlgorithmsFacade arrays;

    function setUp() public {
        voting = new Voting();
        stringReverse = new StringReverse();
        roman = new RomanNumeralsFacade();
        arrays = new ArrayAlgorithmsFacade();
    }

    function test_Voting() public {
        voting.vote("alice");
        voting.vote("alice");
        voting.vote("bob");
        assertEq(voting.getVotes("alice"), 2);
        assertEq(voting.getVotes("bob"), 1);
        voting.resetVotes();
        assertEq(voting.getVotes("alice"), 0);
        assertEq(voting.getVotes("bob"), 0);
    }

    function test_ReverseString() public view {
        assertEq(stringReverse.reverse("abcde"), "edcba");
    }

    function test_IntToRoman() public view {
        assertEq(roman.intToRoman(3), "III");
        assertEq(roman.intToRoman(58), "LVIII");
        assertEq(roman.intToRoman(1994), "MCMXCIV");
    }

    function test_RomanToInt() public view {
        assertEq(roman.romanToInt("III"), 3);
        assertEq(roman.romanToInt("LVIII"), 58);
        assertEq(roman.romanToInt("MCMXCIV"), 1994);
    }

    function test_MergeSorted() public {
        uint256[] memory a = new uint256[](3);
        a[0] = 1;
        a[1] = 3;
        a[2] = 5;
        uint256[] memory b = new uint256[](3);
        b[0] = 2;
        b[1] = 4;
        b[2] = 6;
        uint256[] memory merged = arrays.mergeSorted(a, b);
        assertEq(merged.length, 6);
        assertEq(merged[0], 1);
        assertEq(merged[5], 6);
    }

    function test_BinarySearch() public view {
        uint256[] memory nums = new uint256[](5);
        nums[0] = 1;
        nums[1] = 3;
        nums[2] = 5;
        nums[3] = 7;
        nums[4] = 9;
        assertEq(arrays.binarySearch(nums, 5), 2);
        assertEq(arrays.binarySearch(nums, 11), type(uint256).max);
    }
}
