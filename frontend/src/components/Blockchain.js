import React, { useState, useEffect } from 'react';
import { Box, VStack, Heading, Text, Button, useToast } from '@chakra-ui/react';
import axios from 'axios';

function Blockchain() {
  const [blocks, setBlocks] = useState([]);
  const [loading, setLoading] = useState(false);
  const toast = useToast();

  const fetchBlocks = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/blocks');
      setBlocks(response.data);
    } catch (error) {
      toast({
        title: 'Error fetching blocks',
        description: error.message,
        status: 'error',
        duration: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  const mineBlock = async () => {
    try {
      await axios.post('/mine');
      fetchBlocks();
      toast({
        title: 'Block mined successfully',
        status: 'success',
        duration: 3000,
      });
    } catch (error) {
      toast({
        title: 'Error mining block',
        description: error.message,
        status: 'error',
        duration: 3000,
      });
    }
  };

  useEffect(() => {
    fetchBlocks();
    const interval = setInterval(fetchBlocks, 10000);
    return () => clearInterval(interval);
  }, []);

  return (
    <Box w="100%" bg="white" p={6} borderRadius="lg" shadow="md">
      <VStack spacing={4} align="stretch">
        <Heading size="md">Blockchain</Heading>
        <Button 
          colorScheme="blue" 
          onClick={mineBlock} 
          isLoading={loading}
        >
          Mine New Block
        </Button>
        {blocks.map((block, index) => (
          <Box key={block.Hash} p={4} bg="gray.50" borderRadius="md">
            <Text><strong>Block {index}</strong></Text>
            <Text fontSize="sm">Hash: {block.Hash}</Text>
            <Text fontSize="sm">Previous Hash: {block.PrevBlockHash}</Text>
            <Text fontSize="sm">Timestamp: {new Date(block.Timestamp * 1000).toLocaleString()}</Text>
            <Text fontSize="sm">Transactions: {block.Transactions?.length || 0}</Text>
          </Box>
        ))}
      </VStack>
    </Box>
  );
}

export default Blockchain; 