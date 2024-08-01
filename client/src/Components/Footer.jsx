import React from 'react';
import { Box, Text, Link, Flex, IconButton, VStack, HStack, useColorModeValue } from '@chakra-ui/react';
import { FaGithub, FaTwitter, FaStar, FaHeart, FaEnvelope } from 'react-icons/fa';
import { SiGo } from 'react-icons/si';

function Footer() {
  const bgColor = 'gray.900';
  const textColor = 'gray.300';
  const accentColor = 'teal.400';
  const iconSize = 5; // Chakra UI's size scale

  return (
    <Box bg={bgColor} py={12} color={textColor} borderTop="1px" borderColor="gray.800">
      <Flex direction="column" align="center" maxW="container.lg" mx="auto" px={4}>
        <VStack spacing={6} align="center" w="full">
          <Text fontSize="xl" fontWeight="bold" textAlign="center">
            Made with <Box as={FaHeart} display="inline" color="red.500" mx={1} /> by Ayush Sharma
          </Text>
          
          <HStack spacing={8} justify="center" wrap="wrap">
            <Link href="https://github.com/ayushthe1/streak" isExternal>
              <HStack spacing={2} _hover={{ color: accentColor }}>
                <IconButton
                  aria-label="GitHub"
                  icon={<FaGithub />}
                  variant="ghost"
                  size={iconSize}
                  color="current"
                />
                <Text fontSize="sm">Star on GitHub</Text>
                <Box as={FaStar} color="yellow.400" />
              </HStack>
            </Link>
            
            <Link href="https://x.com/ayushthe5" isExternal>
              <HStack spacing={2} _hover={{ color: accentColor }}>
                <IconButton
                  aria-label="Twitter"
                  icon={<FaTwitter />}
                  variant="ghost"
                  size={iconSize}
                  color="current"
                />
                <Text fontSize="sm">Follow on Twitter</Text>
              </HStack>
            </Link>
          </HStack>
          
          <Link href="mailto:ayushsharmaa101@gmail.com" color={accentColor} _hover={{ textDecoration: 'underline' }} isExternal>
            <HStack spacing={2}>
              <FaEnvelope />
              <Text>ayushsharmaa101@gmail.com</Text>
            </HStack>
          </Link>
          
          <HStack spacing={2} color="gray.500">
            <SiGo />
            <Text fontSize="sm">Powered by Golang</Text>
          </HStack>
        </VStack>
      </Flex>
    </Box>
  );
}

export default Footer;