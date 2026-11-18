// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {BeggingContract} from "../src/BeggingContract.sol";

contract BeggingContractTest is Test {
    BeggingContract begging;
    address owner = makeAddr("owner");
    address donorA = makeAddr("donorA");
    address donorB = makeAddr("donorB");

    function setUp() public {
        vm.prank(owner);
        begging = new BeggingContract();
    }

    function test_DonateAndGetDonation() public {
        vm.deal(donorA, 1 ether);
        vm.prank(donorA);
        begging.donate{value: 0.5 ether}();
        assertEq(begging.getDonation(donorA), 0.5 ether);
        assertEq(address(begging).balance, 0.5 ether);
    }

    function test_WithdrawOnlyOwner() public {
        vm.deal(donorA, 1 ether);
        vm.prank(donorA);
        begging.donate{value: 0.3 ether}();

        vm.prank(donorB);
        vm.expectRevert(BeggingContract.NotOwner.selector);
        begging.withdraw();

        uint256 before = owner.balance;
        vm.prank(owner);
        begging.withdraw();
        assertEq(address(begging).balance, 0);
        assertEq(owner.balance, before + 0.3 ether);
    }

    function test_TopDonors() public {
        vm.deal(donorA, 2 ether);
        vm.deal(donorB, 2 ether);
        vm.prank(donorA);
        begging.donate{value: 1 ether}();
        vm.prank(donorB);
        begging.donate{value: 0.2 ether}();

        (address[] memory leaders, uint256[] memory amounts) = begging.topDonors(3);
        assertEq(leaders[0], donorA);
        assertEq(amounts[0], 1 ether);
    }
}
