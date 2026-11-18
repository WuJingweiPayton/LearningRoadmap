// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

/// @title Array algorithms (homework 01)
library ArrayAlgorithms {
    function mergeSorted(uint256[] memory nums1, uint256[] memory nums2)
        internal
        pure
        returns (uint256[] memory)
    {
        uint256[] memory merged = new uint256[](nums1.length + nums2.length);
        uint256 i;
        uint256 j;
        uint256 k;

        while (i < nums1.length && j < nums2.length) {
            if (nums1[i] <= nums2[j]) {
                merged[k++] = nums1[i++];
            } else {
                merged[k++] = nums2[j++];
            }
        }
        while (i < nums1.length) {
            merged[k++] = nums1[i++];
        }
        while (j < nums2.length) {
            merged[k++] = nums2[j++];
        }
        return merged;
    }

  /// @return index of target, or type(uint256).max if not found
    function binarySearch(uint256[] memory nums, uint256 target) internal pure returns (uint256) {
        uint256 left;
        uint256 right = nums.length;
        while (left < right) {
            uint256 mid = left + (right - left) / 2;
            if (nums[mid] == target) {
                return mid;
            }
            if (nums[mid] < target) {
                left = mid + 1;
            } else {
                right = mid;
            }
        }
        return type(uint256).max;
    }
}

contract ArrayAlgorithmsFacade {
    function mergeSorted(uint256[] calldata nums1, uint256[] calldata nums2)
        external
        pure
        returns (uint256[] memory)
    {
        return ArrayAlgorithms.mergeSorted(nums1, nums2);
    }

    function binarySearch(uint256[] calldata nums, uint256 target) external pure returns (uint256) {
        return ArrayAlgorithms.binarySearch(nums, target);
    }
}
