import React, { useState } from 'react';
import {
  Box,
  VStack,
  Heading,
  FormControl,
  FormLabel,
  Input,
  Button,
  useToast,
} from '@chakra-ui/react';
import axios from 'axios';

function Transaction() {
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');
  const [amount, setAmount] = useState('');
  const [loading, setLoading] = useState(false);
  const toast = useToast();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      setLoading(true);
      await axios.post('/transaction', {
        from,
        to,
        amount: parseFloat(amount),
      });
      toast({
        title: 'Transaction created',
        status: 'success',
        duration: 3000,
      });
      setFrom('');
      setTo('');
      setAmount('');
    } catch (error) {
      toast({
        title: 'Error creating transaction',
        description: error.message,
        status: 'error',
        duration: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box w="100%" bg="white" p={6} borderRadius="lg" shadow="md">
      <VStack spacing={4} align="stretch">
        <Heading size="md">Create Transaction</Heading>
        <form onSubmit={handleSubmit}>
          <VStack spacing={4}>
            <FormControl>
              <FormLabel>From Address</FormLabel>
              <Input
                value={from}
                onChange={(e) => setFrom(e.target.value)}
                placeholder="Sender's address"
              />
            </FormControl>
            <FormControl>
              <FormLabel>To Address</FormLabel>
              <Input
                value={to}
                onChange={(e) => setTo(e.target.value)}
                placeholder="Recipient's address"
              />
            </FormControl>
            <FormControl>
              <FormLabel>Amount</FormLabel>
              <Input
                type="number"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="Amount to send"
              />
            </FormControl>
            <Button
              type="submit"
              colorScheme="green"
              isLoading={loading}
              w="100%"
            >
              Send Transaction
            </Button>
          </VStack>
        </form>
      </VStack>
    </Box>
  );
}

export default Transaction; 