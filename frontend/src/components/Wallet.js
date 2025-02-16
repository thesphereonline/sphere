import React, { useState, useEffect } from 'react';
import {
  Box,
  VStack,
  Heading,
  Text,
  Button,
  useToast,
  HStack,
  Input,
  FormControl,
  FormLabel,
} from '@chakra-ui/react';
import axios from 'axios';

function Wallet() {
  const [wallet, setWallet] = useState(null);
  const [balance, setBalance] = useState(0);
  const [loading, setLoading] = useState(false);
  const [importAddress, setImportAddress] = useState('');
  const toast = useToast();

  const fetchWalletInfo = async () => {
    try {
      const response = await axios.get('/status');
      setWallet({
        address: response.data.minerAddress,
      });
      
      // Fetch balance for the wallet
      const balanceResponse = await axios.get(`/balance/${response.data.minerAddress}`);
      setBalance(balanceResponse.data.balance);
    } catch (error) {
      toast({
        title: 'Error fetching wallet info',
        description: error.message,
        status: 'error',
        duration: 3000,
      });
    }
  };

  const createNewWallet = async () => {
    try {
      setLoading(true);
      const response = await axios.post('/wallet/create');
      setWallet(response.data);
      toast({
        title: 'New wallet created',
        status: 'success',
        duration: 3000,
      });
    } catch (error) {
      toast({
        title: 'Error creating wallet',
        description: error.message,
        status: 'error',
        duration: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  const importWallet = async () => {
    try {
      setLoading(true);
      const response = await axios.post('/wallet/import', { address: importAddress });
      setWallet(response.data);
      setImportAddress('');
      toast({
        title: 'Wallet imported successfully',
        status: 'success',
        duration: 3000,
      });
    } catch (error) {
      toast({
        title: 'Error importing wallet',
        description: error.message,
        status: 'error',
        duration: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchWalletInfo();
    const interval = setInterval(fetchWalletInfo, 10000);
    return () => clearInterval(interval);
  }, []);

  return (
    <Box w="100%" bg="white" p={6} borderRadius="lg" shadow="md">
      <VStack spacing={4} align="stretch">
        <Heading size="md">Wallet</Heading>
        
        {wallet ? (
          <VStack align="stretch" spacing={2}>
            <Text><strong>Address:</strong> {wallet.address}</Text>
            <Text><strong>Balance:</strong> {balance} coins</Text>
            <Button
              colorScheme="blue"
              onClick={fetchWalletInfo}
              isLoading={loading}
            >
              Refresh Balance
            </Button>
          </VStack>
        ) : (
          <Text>No wallet loaded</Text>
        )}

        <Box borderTop="1px" borderColor="gray.200" pt={4}>
          <VStack spacing={4}>
            <Button
              colorScheme="green"
              onClick={createNewWallet}
              isLoading={loading}
              w="100%"
            >
              Create New Wallet
            </Button>

            <FormControl>
              <FormLabel>Import Existing Wallet</FormLabel>
              <HStack>
                <Input
                  value={importAddress}
                  onChange={(e) => setImportAddress(e.target.value)}
                  placeholder="Enter wallet address"
                />
                <Button
                  colorScheme="blue"
                  onClick={importWallet}
                  isLoading={loading}
                  isDisabled={!importAddress}
                >
                  Import
                </Button>
              </HStack>
            </FormControl>
          </VStack>
        </Box>
      </VStack>
    </Box>
  );
}

export default Wallet; 