import React, { Component } from 'react';
import axios from 'axios';
import {
  Container,
  FormControl,
  FormLabel,
  FormErrorMessage,
  Box,
  Input,
  Stack,
  Button,
  Heading,
  Text,
  VStack,
  InputGroup,
  InputRightElement,
  Icon,
} from '@chakra-ui/react';
import { EditIcon, ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import { FaUserPlus } from 'react-icons/fa';
import { Navigate } from 'react-router-dom';

class Register extends Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      message: '',
      isInvalid: false,
      endpoint: 'https://api.ayushsharma.co.in/api/register',
      redirect: false,
      redirectTo: '/chat?u=',
      showPassword: false,
    };
  }

  onChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };

  onSubmit = async e => {
    e.preventDefault();

    try {
      const res = await axios.post(this.state.endpoint, {
        username: this.state.username,
        password: this.state.password,
      }, { withCredentials: true });

      console.log('register', res);
      if (res.status) {
        const redirectTo = this.state.redirectTo + this.state.username;
        this.setState({ redirect: true, redirectTo });
      } else {
        this.setState({ message: res.data.message, isInvalid: true });
      }
    } catch (error) {
      console.log(error);
      this.setState({ message: 'Something went wrong', isInvalid: true });
    }
  };

  togglePasswordVisibility = () => {
    this.setState(prevState => ({ showPassword: !prevState.showPassword }));
  };

  render() {
    if (this.state.redirect) {
      return <Navigate to={this.state.redirectTo} replace={true} />;
    }

    return (
      <Box bg="#2f3349" minH="100vh" py={20}>
        <Container maxW="md">
          <VStack spacing={8} align="stretch">
            <VStack spacing={2} align="center">
              <Icon as={FaUserPlus} w={10} h={10} color="purple.400" />
              <Heading color="white">Create Your Account</Heading>
              <Text color="gray.400">Join Streak and start chatting today!</Text>
            </VStack>
            <Box bg="gray.800" borderRadius="xl" p={8} boxShadow="xl">
              <form onSubmit={this.onSubmit}>
                <Stack spacing={6}>
                  <FormControl isInvalid={this.state.isInvalid}>
                    <FormLabel color="gray.300">Username</FormLabel>
                    <Input
                      type="text"
                      placeholder="Enter your username (no special char)"
                      name="username"
                      value={this.state.username}
                      onChange={this.onChange}
                      bg="gray.700"
                      color="white"
                      border="none"
                      _placeholder={{ color: 'gray.400' }}
                      _focus={{ borderColor: 'purple.400', boxShadow: '0 0 0 1px #9F7AEA' }}
                    />
                    {this.state.isInvalid && (
                      <FormErrorMessage>{this.state.message}</FormErrorMessage>
                    )}
                  </FormControl>
                  <FormControl>
                    <FormLabel color="gray.300">Password</FormLabel>
                    <InputGroup>
                      <Input
                        type={this.state.showPassword ? 'text' : 'password'}
                        placeholder="Enter your password"
                        name="password"
                        value={this.state.password}
                        onChange={this.onChange}
                        bg="gray.700"
                        color="white"
                        border="none"
                        _placeholder={{ color: 'gray.400' }}
                        _focus={{ borderColor: 'purple.400', boxShadow: '0 0 0 1px #9F7AEA' }}
                      />
                      <InputRightElement width="3rem">
                        <Button h="1.5rem" size="sm" onClick={this.togglePasswordVisibility} bg="transparent">
                          {this.state.showPassword ? (
                            <ViewOffIcon color="gray.400" />
                          ) : (
                            <ViewIcon color="gray.400" />
                          )}
                        </Button>
                      </InputRightElement>
                    </InputGroup>
                  </FormControl>
                  <Button
                    leftIcon={<EditIcon />}
                    colorScheme="purple"
                    variant="solid"
                    type="submit"
                    size="lg"
                  >
                    Register
                  </Button>
                </Stack>
              </form>
            </Box>
          </VStack>
        </Container>
      </Box>
    );
  }
}

export default Register;