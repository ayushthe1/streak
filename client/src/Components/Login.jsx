// import React, { Component } from 'react';
// import axios from 'axios';

// import {
//   Container,
//   FormControl,
//   FormLabel,
//   FormErrorMessage,
//   Text,
//   Box,
//   Input,
//   Stack,
//   Button,
// } from '@chakra-ui/react';

// import { Navigate } from 'react-router-dom';
// import { EditIcon } from '@chakra-ui/icons';

// class Login extends Component {
//   constructor(props) {
//     super(props);
//     this.state = {
//       username: '',
//       password: '',
//       message: '',
//       isInvalid: false,
//       endpoint: 'https://api.ayushsharma.co.in/api/login',
//       redirect: false,
//       redirectTo: '/chat?u=',
//     };
//   }

//   // on change of input, set the value to the message state
//   onChange = event => {
//     this.setState({ [event.target.name]: event.target.value });
//   };

//   onSubmit = async e => {
//     e.preventDefault();

//     try {
//       const res = await axios.post(this.state.endpoint, {
//         username: this.state.username,
//         password: this.state.password,
//       });

//       console.log('register', res);
//       if (res.status) {
//         const redirectTo = this.state.redirectTo + this.state.username;
//         this.setState({ redirect: true, redirectTo });
//       } else {
//         // on failed
//         this.setState({ message: res.data.message, isInvalid: true });
//       }
//     } catch (error) {
//       console.log(error);
//       this.setState({ message: 'something went wrong', isInvalid: true });
//     }
//   };

//   render() {
//     return (
//       <div>
//         {this.state.redirect && (
//           <Navigate to={this.state.redirectTo} replace={true}></Navigate>
//         )}

//         <Container marginBlockStart={10} textAlign={'left'} maxW="2xl">
//           <Box borderRadius="lg" padding={10} borderWidth="2px">
//             <Stack spacing={5}>
//               <FormControl isInvalid={this.state.isInvalid}>
//                 <FormLabel>Username</FormLabel>
//                 <Input
//                   type="text"
//                   placeholder="Username"
//                   name="username"
//                   value={this.state.username}
//                   onChange={this.onChange}
//                 />
//               </FormControl>
//               <FormControl isInvalid={this.state.isInvalid}>
//                 <FormLabel>Password</FormLabel>
//                 <Input
//                   type="password"
//                   placeholder="Password"
//                   name="password"
//                   value={this.state.password}
//                   onChange={this.onChange}
//                 />
//                 {!this.state.isInvalid ? (
//                   ''
//                 ) : (
//                   <FormErrorMessage>
//                     invalid username or password
//                   </FormErrorMessage>
//                 )}
//               </FormControl>
//               <Button
//                 size="lg"
//                 leftIcon={<EditIcon />}
//                 colorScheme="green"
//                 variant="solid"
//                 type="submit"
//                 onClick={this.onSubmit}
//               >
//                 Login
//               </Button>
//             </Stack>
//             <Box paddingTop={3}>
//               <Text as="i" fontSize={'lg'} color={'red'}>
//                 {this.state.message}
//               </Text>
//             </Box>
//           </Box>
//         </Container>
//       </div>
//     );
//   }
// }

// export default Login;

import React, { Component } from 'react';
import axios from 'axios';

import {
  Container,
  FormControl,
  FormLabel,
  FormErrorMessage,
  Text,
  Box,
  Input,
  Stack,
  Button,
} from '@chakra-ui/react';

import { Navigate } from 'react-router-dom';
import { EditIcon } from '@chakra-ui/icons';

class Login extends Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      message: '',
      isInvalid: false,
      endpoint: 'https://api.ayushsharma.co.in/api/login',
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

      console.log('login response ', res);
      if (res.status === 200) {
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
      <Container
      marginBlockStart={{ base: 5, md: 10 }}
      textAlign="center"
      maxW={{ base: '90%', sm: 'sm' }}
      padding={4}
    >
      <Box
        borderRadius="lg"
        p={{ base: 6, md: 10 }}
        borderWidth="2px"
        bg="gray.800"
        color="white"
        boxShadow="lg"
      >
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
                bg="gray.700"
                color="white"
                border="none"
                _placeholder={{ color: 'gray.400' }}
              />
            </FormControl>
            <FormControl isInvalid={this.state.isInvalid}>
              <FormLabel>Password</FormLabel>
              <Input
                type="password"
                placeholder="Password"
                name="password"
                value={this.state.password}
                onChange={this.onChange}
                bg="gray.700"
                color="white"
                border="none"
                _placeholder={{ color: 'gray.400' }}
              />
              {this.state.isInvalid && (
                <FormErrorMessage>
                  Invalid username or password
                </FormErrorMessage>
              )}
            </FormControl>
            <Button
              size="lg"
              leftIcon={<EditIcon />}
              colorScheme="purple"
              variant="solid"
              type="submit"
              boxShadow="md"
              _hover={{ boxShadow: 'xl' }}
              transition="box-shadow 0.2s"
            >
              Login
            </Button>
          </Stack>
        </form>
        <Box pt={3}>
          <Text as="i" fontSize="lg" color="red.300">
            {this.state.message}
          </Text>
        </Box>
      </Box>
    </Container>
    );
    
  }
}

export default Login;

