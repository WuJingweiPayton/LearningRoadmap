// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

/// @title Begging / donation contract (homework 02)
contract BeggingContract {
    address public immutable owner;
    mapping(address => uint256) public donations;

    address[] private _donors;
    mapping(address => bool) private _isDonor;

    event Donation(address indexed donor, uint256 amount);

    error NotOwner();
    modifier onlyOwner() {
        if (msg.sender != owner) revert NotOwner();
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    function donate() external payable {
        _recordDonation();
    }

    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        if (balance == 0) {
            return;
        }
        payable(owner).transfer(balance);
    }

    function getDonation(address donor) external view returns (uint256) {
        return donations[donor];
    }

    function topDonors(uint256 count) external view returns (address[] memory leaders, uint256[] memory amounts) {
        uint256 n = count > _donors.length ? _donors.length : count;
        if (n == 0) {
            return (new address[](0), new uint256[](0));
        }
        if (n > 3) {
            n = 3;
        }

        address[] memory addrs = new address[](_donors.length);
        uint256[] memory vals = new uint256[](_donors.length);
        for (uint256 i = 0; i < _donors.length; i++) {
            addrs[i] = _donors[i];
            vals[i] = donations[_donors[i]];
        }

        for (uint256 i = 0; i < _donors.length; i++) {
            for (uint256 j = i + 1; j < _donors.length; j++) {
                if (vals[j] > vals[i]) {
                    (vals[i], vals[j]) = (vals[j], vals[i]);
                    (addrs[i], addrs[j]) = (addrs[j], addrs[i]);
                }
            }
        }

        leaders = new address[](n);
        amounts = new uint256[](n);
        for (uint256 i = 0; i < n; i++) {
            leaders[i] = addrs[i];
            amounts[i] = vals[i];
        }
    }

    receive() external payable {
        _recordDonation();
    }

    function _recordDonation() private {
        if (msg.value == 0) {
            return;
        }
        donations[msg.sender] += msg.value;
        if (!_isDonor[msg.sender]) {
            _isDonor[msg.sender] = true;
            _donors.push(msg.sender);
        }
        emit Donation(msg.sender, msg.value);
    }
}
