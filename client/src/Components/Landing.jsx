import React from 'react';
import { Link } from 'react-router-dom';
import { 
  Box, Button, Container, Flex, Grid, Heading, Text, Stack, 
  Icon, VStack, HStack, useColorModeValue, Image
} from '@chakra-ui/react';
import { FaRobot, FaSearch, FaFileAlt, FaUserFriends } from 'react-icons/fa';
import { ArrowForwardIcon } from '@chakra-ui/icons';

const Feature = ({ title, text, icon }) => {
  return (
    <VStack 
      align="center" 
      spacing={4} 
      bg="whiteAlpha.100" 
      p={6} 
      borderRadius="lg" 
      transition="all 0.3s"
      _hover={{ transform: 'translateY(-5px)', boxShadow: 'lg' }}
    >
      <Icon as={icon} w={10} h={10} color="purple.400" />
      <Text fontWeight="bold" fontSize="xl" textAlign="center">{title}</Text>
      <Text textAlign="center">{text}</Text>
    </VStack>
  );
};

function Landing() {
  const bgColor = '#2f3349';
  const textColor = 'white';
  const accentColor = 'purple.400';

  return (
    <Box>
    <Box 
  backgroundImage="url('https://images.unsplash.com/photo-1551434678-e076c223a692?ixlib=rb-1.2.1&auto=format&fit=crop&w=2850&q=80')"
  backgroundPosition="center"
  backgroundRepeat="no-repeat"
  backgroundSize="cover"
>
  <Box bg="rgba(47, 51, 73, 0.85)" py={20} minHeight="80vh" display="flex" alignItems="center">
    <Container maxW="container.xl">
      <VStack spacing={8} align="center" textAlign="center">
        <Heading as="h1" size="3xl" lineHeight="shorter">
          Welcome to <Text as="span" color={accentColor}>Streak</Text>
        </Heading>
        <Text fontSize="xl" maxW="2xl">
          Experience the future of communication with our AI-powered chat platform. 
          Connect, share, and discover with unprecedented ease and intelligence.
        </Text>
        <Stack direction={{ base: "column", sm: "row" }} spacing={4} pt={4}>
          <Button 
            as={Link} 
            to="/register" 
            size="lg" 
            colorScheme="purple" 
            rightIcon={<ArrowForwardIcon />}
          >
            Register
          </Button>
          <Button 
            as={Link} 
            to="/login" 
            size="lg" 
            variant="outline" 
            colorScheme="purple"
          >
            Login
          </Button>
        </Stack>
      </VStack>
    </Container>
  </Box>
  </Box>

      {/* Features Section */}
      <Box py={20}>
        <Container maxW="container.xl">
          <VStack spacing={16}>
            <Heading as="h2" size="xl" textAlign="center">
              Why Choose Streak?
            </Heading>
            <Grid templateColumns={{ base: "1fr", md: "1fr 1fr" }} gap={8}>
              <Feature 
                icon={FaUserFriends}
                title="Real-time Chat"
                text="Connect instantly with friends and colleagues through our lightning-fast chat system."
              />
              <Feature 
                icon={FaFileAlt}
                title="File Sharing"
                text="Share files effortlessly during your conversations, enhancing collaboration."
              />
              <Feature 
                icon={FaRobot}
                title="AI Chatbot"
                text="Get instant answers and assistance from our intelligent AI chatbot companion."
              />
              <Feature 
                icon={FaSearch}
                title="Semantic Search"
                text="Find past conversations and information with our powerful semantic search feature."
              />
            </Grid>
          </VStack>
        </Container>
      </Box>

      {/* Call to Action Section */}
      <Box bg="whiteAlpha.100" py={20}>
        <Container maxW="container.xl">
          <VStack spacing={8} textAlign="center">
            <Heading as="h2" size="xl">
              Ready to Transform Your Communication?
            </Heading>
            <Text fontSize="lg" maxW="2xl">
              Join thousands of users already experiencing the power of Streak.
              Start your journey to smarter, more efficient communication today.
            </Text>
            <Button 
              as={Link} 
              to="/register" 
              size="lg" 
              colorScheme="purple" 
              rightIcon={<ArrowForwardIcon />}
            >
              Start Your Streak Today
            </Button>
          </VStack>
        </Container>
      </Box>
    </Box>
  );
}

export default Landing;