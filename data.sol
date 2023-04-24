// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyNFT is ERC721, Pausable, Ownable {
    uint256 public constant MAX_NFTS = 5;
    uint256 public constant MINT_PRICE = 0.1 ether;
    uint256 public constant MINT_DURATION_START = 1641532800; // 7 Jan 2022 00:00:00 GMT
    uint256 public constant MINT_DURATION_END = 1642144000; // 14 Jan 2022 00:00:00 GMT
    mapping(address => mapping(string => bool)) private _receipts;

    constructor() ERC721("MyNFT", "MNFT") {}

    function mintNFT(string memory receipt, string memory name, string memory description, string memory image) public payable whenNotPaused {
        require(block.timestamp >= MINT_DURATION_START && block.timestamp < MINT_DURATION_END, "Minting period has ended");
        require(msg.value >= MINT_PRICE, "Insufficient payment");

        require(!_receipts[msg.sender][receipt], "This receipt has already been used");
        require(totalSupply() < MAX_NFTS, "All NFTs have been minted");

        uint256 tokenId = totalSupply() + 1;
        _mint(msg.sender, tokenId);
        _setTokenURI(tokenId, image);
        _receipts[msg.sender][receipt] = true;

        emit MintedNFT(msg.sender, tokenId, name, description, image);
    }

    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        payable(msg.sender).transfer(balance);
    }

    event MintedNFT(address indexed owner, uint256 indexed tokenId, string name, string description, string image);
}
