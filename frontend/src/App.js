import React from 'react';
import { ChakraProvider, Box, VStack, Heading, Container } from '@chakra-ui/react';
import Blockchain from './components/Blockchain';
import Transaction from './components/Transaction';
import Wallet from './components/Wallet';
import Peers from './components/Peers';

function App() {
  return (
    <ChakraProvider>
      <Box minH="100vh" bg="gray.100" py={8}>
        <Container maxW="container.xl">
          <VStack spacing={8}>
            <Heading>Blockchain Network Dashboard</Heading>
            <Wallet />
            <Transaction />
            <Blockchain />
            <Peers />
          </VStack>
        </Container>
      </Box>
    </ChakraProvider>
  );
}

export default App; 