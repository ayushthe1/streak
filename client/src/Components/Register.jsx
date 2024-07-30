import React, { Component } from 'react';
import axios from 'axios';

import {
  Container,
  FormControl,
  FormLabel,
  FormErrorMessage,
  FormHelperText,
  Box,
  Input,
  Stack,
  Button,
} from '@chakra-ui/react';

import { EditIcon } from '@chakra-ui/icons';
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
    };
  }

  // on change of input, set the value to the message state
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
        // on failed
        this.setState({ message: res.data.message, isInvalid: true });
      }
    } catch (error) {
      console.log(error);
      this.setState({ message: 'something went wrong', isInvalid: true });
    }
  };

  render() {
    if (this.state.redirect) {
      return <Navigate to={this.state.redirectTo} replace={true} />;
    }

    return (
      <Container marginBlockStart={10} textAlign={'left'} maxW="2xl">
        <Box borderRadius="lg" padding={10} borderWidth="2px">
          <form onSubmit={this.onSubmit}>
            <Stack spacing={5}>
              <FormControl isInvalid={this.state.isInvalid}>
                <FormLabel>Username</FormLabel>
                <Input
                  type="text"
                  placeholder="Username"
                  name="username"
                  value={this.state.username}
                  onChange={this.onChange}
                />
                {this.state.isInvalid && (
                  <FormErrorMessage>{this.state.message}</FormErrorMessage>
                )}
              </FormControl>
              <FormControl>
                <FormLabel>Password</FormLabel>
                <Input
                  type="password"
                  placeholder="Password"
                  name="password"
                  value={this.state.password}
                  onChange={this.onChange}
                />
                <FormHelperText>use a dummy password</FormHelperText>
              </FormControl>
              <Button
                size="lg"
                leftIcon={<EditIcon />}
                colorScheme="green"
                variant="solid"
                type="submit"
              >
                Register
              </Button>
            </Stack>
          </form>
        </Box>
      </Container>
    );
  }
}

export default Register;
