// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.6.8;
pragma experimental ABIEncoderV2;

contract Staking {
    event Delegated(address delegator, string validator, uint256 amount);
    event Undelegated(address delegator, string validator, uint256 amount);
    event Redelegated(
        address delegator,
        string validatorSrc,
        string validatorDest,
        uint256 amount
    );
    event Withdrew(address delegator, string validator);

    function delegate(string calldata validator, uint256 amount) external {
        emit Delegated(msg.sender, validator, amount);
    }

    function undelegate(string calldata validator, uint256 amount) external {
        emit Undelegated(msg.sender, validator, amount);
    }

    function redelegate(
        string calldata validatorSrc,
        string calldata validatorDest,
        uint256 amount
    ) external {
        emit Redelegated(msg.sender, validatorSrc, validatorDest, amount);
    }

    function withdraw(string calldata validator) external {
        emit Withdrew(msg.sender, validator);
    }
}
