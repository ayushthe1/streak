import React, { useState } from 'react';
import {
  Box,
  Input,
  Button,
  InputGroup,
  InputRightElement,
  useToast,
} from '@chakra-ui/react';
import { SearchIcon } from '@chakra-ui/icons';
import axios from 'axios';

const SearchBar = ({ from, onSearchResults }) => {
  const [query, setQuery] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const toast = useToast();

  const handleSearch = async () => {
    if (!query.trim()) {
      toast({
        title: 'Error',
        description: 'Please enter a search query',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    setIsLoading(true);
    try {
      const response = await axios.post('https://api.ayushsharma.co.in/api/wv', { query, from });
      onSearchResults(response.data);
    } catch (error) {
      console.error('Error fetching results:', error);
      toast({
        title: 'Error',
        description: 'Failed to fetch results. Please try again.',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box mb={4}>
      <InputGroup size="md">
        <Input
          pr="4.5rem"
          placeholder="Search..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <InputRightElement width="4.5rem">
          <Button
            h="1.75rem"
            size="sm"
            onClick={handleSearch}
            isLoading={isLoading}
          >
            <SearchIcon />
          </Button>
        </InputRightElement>
      </InputGroup>
    </Box>
  );
};

export default SearchBar;