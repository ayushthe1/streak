import React from 'react';
import { Link } from 'react-router-dom';
import { Box, Button, Container, Stack, Text, useColorModeValue } from '@chakra-ui/react';
import { EditIcon, ArrowForwardIcon } from '@chakra-ui/icons';

function Landing() {
  const bgColor = useColorModeValue('gray.800', 'gray.800');
  const textColor = useColorModeValue('white', 'white');
  const buttonColorScheme = 'purple';

  return (
    <Container maxW="2xl" mt="3rem" centerContent>
      <Box padding="5" bg={bgColor} borderRadius="lg" textAlign="center" color={textColor} boxShadow="lg">
        <Text fontSize="4xl" fontWeight="bold" mb={5}>
          Welcome to Streak!
        </Text>
        <Stack direction="row" spacing={7} justify="center">
          <Link to="register">
            <Button
              size="lg"
              leftIcon={<EditIcon />}
              colorScheme={buttonColorScheme}
              variant="solid"
            >
              Register
            </Button>
          </Link>
          <Link to="login">
            <Button
              size="lg"
              rightIcon={<ArrowForwardIcon />}
              colorScheme={buttonColorScheme}
              variant="outline"
            >
              Login
            </Button>
          </Link>
        </Stack>
      </Box>
    </Container>
  );
}

export default Landing;
