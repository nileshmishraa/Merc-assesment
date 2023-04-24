This assessment is about developing an NFT portal in order to assess the technical skills of the  candidate in the areas of Web3.0, Full stack and API.  
The Assessment consists of 3 parts :
1. Smart contracts :  
   a. Develop and Deploy NFT smart contracts  
   b. The smart contracts should be able to mint NFT  
   i. Mint only valid certain duration (example between 7 Jan to 14 Jan 2023)  
   ii. Mint only once for each wallet and Receipt (refer to 3.i)
   iii. The receipt will have to be store in smart contract state
   iv. Only able to mint 5 NFT  
   v. The NFT should have metadata (name, description, image)  
   c. Script to deploy the smart contract
2. WebApp :  
   a. React app with any preferred React framework can be used  
   b. Web3 integration with web3.js or ether.js  
   c. Collect user input e.g. NRIC
   d. Interact with Smart contract by Claim (mint) NFT with the connected wallet and  Receipt (refer to 3.i)
   e. The App should display the NFT image from NFT metadata  
   f. The necessary error handlings to be developed
3. API:
   a. Golang API
   b. Any preferred framework (example gin-gonic)
   c. The API will collect National Registration Identity Card (NRIC) and wallet address  from WebApp
   d. POST API body: NRIC and wallet address
   e. NRIC must be unique
   f. Wallet address can only be associated with 1 NRIC
   g. Store into RDBS (PostgreSQL, MySQL, etc) for the unique NRIC and wallet  
   address
   h. Provide the docker-compose.yaml script for the RDBS stack
   i. POST API Response with a Receipt produce by encrypt or hash the API body, you  would need to explain why you choose one mechanism over the other (encrypt vs  
   hash) 
