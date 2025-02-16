import React, { useState, useEffect } from 'react';
import {
  Box,
  VStack,
  Heading,
  Text,
  Input,
  Button,
  List,
  ListItem,
  useToast,
} from '@chakra-ui/react';
import axios from 'axios';

function Peers() {
  const [peers, setPeers] = useState([]);
  const [newPeer, setNewPeer] = useState('');
  const [loading, setLoading] = useState(false);
  const toast = useToast();

  const fetchPeers = async () => {
    try {
      const response = await axios.get('/peers');
      setPeers(response.data);
    } catch (error) {
      toast({
        title: 'Error fetching peers',
        description: error.message,
        status: 'error',
        duration: 3000,
      });
    }
  };

  const addPeer = async () => {
    try {
      setLoading(true);
      await axios.post('/peers', { address: newPeer });
      setNewPeer('');
      fetchPeers();
      toast({
        title: 'Peer added successfully',
        status: 'success',
        duration: 3000,
      });
    } catch (error) {
      toast({
        title: 'Error adding peer',
        description: error.message,
        status: 'error',
        duration: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPeers();
    const interval = setInterval(fetchPeers, 10000);
    return () => clearInterval(interval);
  }, []);

  return (
    <Box w="100%" bg="white" p={6} borderRadius="lg" shadow="md">
      <VStack spacing={4} align="stretch">
        <Heading size="md">Network Peers</Heading>
        <Box>
          <Input
            value={newPeer}
            onChange={(e) => setNewPeer(e.target.value)}
            placeholder="Enter peer address (e.g., localhost:8081)"
            mb={2}
          />
          <Button
            colorScheme="blue"
            onClick={addPeer}
            isLoading={loading}
            w="100%"
          >
            Add Peer
          </Button>
        </Box>
        <List spacing={2}>
          {peers.map((peer) => (
            <ListItem key={peer} p={2} bg="gray.50" borderRadius="md">
              <Text>{peer}</Text>
            </ListItem>
          ))}
        </List>
      </VStack>
    </Box>
  );
}

export default Peers; 