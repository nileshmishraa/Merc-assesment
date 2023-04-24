import React, { useState } from 'react';
import Web3 from 'web3';
import MyNFTContract from './contracts/MyNFT.json';

// running Ganache on 8545
const web3 = new Web3(Web3.givenProvider || 'http://localhost:8545');

function App() {
  const [nric, setNric] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [nft, setNft] = useState(null);

  const handleNricChange = (event) => {
    setNric(event.target.value);
  };

  const handleClaim = async () => {
    setLoading(true);

    try {
      // Load the smart contract instance
      const networkId = await web3.eth.net.getId();
      const deployedNetwork = MyNFTContract.networks[networkId];
      if (!deployedNetwork) {
        throw new Error(`Contract not deployed to network with id ${networkId}`);
      }
      const contract = new web3.eth.Contract(MyNFTContract.abi, deployedNetwork.address);

      // Mint a new NFT and store the user's NRIC as the receipt
      const accounts = await web3.eth.getAccounts();
      if (!accounts.length) {
        throw new Error('No accounts found. Please connect to an Ethereum wallet.');
      }
      const result = await contract.methods.claim(nric).send({ from: accounts[0] });
      if (!result.events.Transfer || !result.events.Transfer.returnValues) {
        throw new Error('Failed to retrieve NFT metadata. Please check your transaction history.');
      }
      const tokenId = result.events.Transfer.returnValues.tokenId;

      // Get the metadata of the newly minted NFT
      const metadata = await contract.methods.tokenURI(tokenId).call();

      setNft({ tokenId, metadata });
      setError(null);
    } catch (err) {
      console.error(err);
      setError(err.message || 'Failed to claim NFT');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App">
      <h1>MyNFT Claim Form</h1>
      <label>
        NRIC:
        <input type="text" value={nric} onChange={handleNricChange} />
      </label>
      <button onClick={handleClaim} disabled={!nric || loading}>
        {loading ? 'Loading...' : 'Claim NFT'}
      </button>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {nft && (
        <>
          <p>Token ID: {nft.tokenId}</p>
          <img src={nft.metadata.image} alt={nft.metadata.name} />
          <p>{nft.metadata.name}</p>
          <p>{nft.metadata.description}</p>
        </>
      )}
    </div>
  );
}

export default App;
