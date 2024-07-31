import React from 'react';
import { Box, Text, Link, Flex, IconButton, Center, VStack, HStack, useColorModeValue } from '@chakra-ui/react';
import { FaGithub, FaTwitter, FaStar } from 'react-icons/fa';
import { MdFavorite } from 'react-icons/md';

function Footer() {
  const iconColor = useColorModeValue('gray.300', 'gray.500');
  const iconSize = '24px'; // Icon size
  const bgColor = useColorModeValue('gray.900', 'gray.800');

  return (
    <Box bg={bgColor} py={8} color="white" borderTop="1px" borderColor="gray.700">
      <Center>
        <VStack spacing={6}>
          <Text fontSize="xl" fontWeight="bold">
            Made with <Box as={MdFavorite} color="red.500" display="inline" boxSize={iconSize} /> by Ayush Sharma
          </Text>
          <HStack spacing={6}>
            <Link href="https://github.com/ayushthe1/streak" isExternal>
              <IconButton
                aria-label="GitHub"
                icon={<FaGithub size={iconSize} />}
                variant="ghost"
                color={iconColor}
                _hover={{ color: 'white' }}
              />
              <HStack spacing={1}>
                <Box as={FaStar} color="yellow.400" boxSize={iconSize} />
                <Text fontSize="sm">Star the project</Text>
              </HStack>
            </Link>
            <Link href="https://x.com/ayushthe5" isExternal>
              <IconButton
                aria-label="Twitter"
                icon={<FaTwitter size={iconSize} />}
                variant="ghost"
                color={iconColor}
                _hover={{ color: 'blue.400' }}
              />
            </Link>
          </HStack>
          <Link href="mailto:ayushsharmaa101@gmail.com" color="teal.400" isExternal>
            your-email@example.com
          </Link>
          <Text fontSize="sm" color="gray.500">
            Powered by Golang
          </Text>
        </VStack>
      </Center>
    </Box>
  );
}

export default Footer;
