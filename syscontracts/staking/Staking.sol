//SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

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

    function delegate(string memory validator, uint256 amount) external {
        emit Delegated(msg.sender, validator, amount);
    }

    function undelegate(string memory validator, uint256 amount) external {
        emit Undelegated(msg.sender, validator, amount);
    }

    function redelegate(
        string memory validatorSrc,
        string memory validatorDest,
        uint256 amount
    ) external {
        emit Redelegated(msg.sender, validatorSrc, validatorDest, amount);
    }

    function withdraw(string memory validator) external {
        emit Withdrew(msg.sender, validator);
    }
}
