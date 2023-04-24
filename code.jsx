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
    // validate the NRIC entered by the user
    if (!nric || !/^[STFG]\d{7}[A-Z]$/.test(nric)) {
      setError('Please enter a valid NRIC.');
      return;
    }

    try {
      setLoading(true);

      // Load the smart contract instance
      const networkId = await web3.eth.net.getId();
      const deployedNetwork = MyNFTContract.networks[networkId];
      if (!deployedNetwork) {
        setError('MyNFT contract is not deployed on the current network.');
        return;
      }
      const contract = new web3.eth.Contract(MyNFTContract.abi, deployedNetwork.address);

      // Mint a new NFT and store the user's NRIC as the receipt
      const defaultAccount = await web3.eth.getAccounts().then((accounts) => accounts[0]);
      const result = await contract.methods.claim(nric).send({ from: defaultAccount });

      // check for any transaction errors
      if (result.status !== '0x1') {
        setError('Failed to claim NFT. Please try again later.');
        return;
      }

      // Get the metadata of the newly minted NFT
      const tokenId = await contract.methods.getTokenByIndex(0).call();
      const metadata = await contract.methods.tokenURI(tokenId).call();

      setNft({ tokenId, metadata });
      setLoading(false);
      setError(null);
    } catch (error) {
      console.error(error);
      setLoading(false);
      setError('Something went wrong, please try again later.');
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
