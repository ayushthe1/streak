
import React, { Component } from 'react';
import axios from 'axios';

import SocketConnection from '../../socket-connection';

import {
  Container,
  Flex,
  Textarea,
  Box,
  FormControl,
  FormErrorMessage,
  InputGroup,
  InputRightElement,
  Button,
  Input,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  VStack,
  Text,
} from '@chakra-ui/react';

import ChatHistory from './ChatHistory';
import ContactList from './ContactList';

class Chat extends Component {
  constructor(props) {
    super(props);
    this.state = {
      socketConn: '',
      username: '',
      message: '',
      to: '',
      isInvalid: false,
      endpoint: 'http://localhost:3000/api',
      contact: '',
      contacts: [],
      renderContactList: [],
      chats: [],
      chatHistory: [],
      activities: [],
      msgs: [],
      file: null,
      fileUrl: ''
    };
  }

  componentDidMount = async () => {
    const queryParams = new URLSearchParams(window.location.search);
    const user = queryParams.get('u');
    this.setState({ username: user });
    this.getContacts(user); // get all contacts of the user

    const conn = new SocketConnection();
    await this.setState({ socketConn: conn });
    // conn.connect(msg => console.log('message received'));
    // connect to ws connection
    this.state.socketConn.connect(message => {
      const msg = JSON.parse(message.data);
      console.log("Message is :", msg);

      if (msg.type === 'activity') {
        this.setState((prevState) => ({
          activities: [msg, ...prevState.activities],
        }));
      } else {

      // update UI only when message is between from and to
      if (this.state.username === msg.to || this.state.username === msg.from) {
        this.setState(
          {
            chats: [...this.state.chats, msg],
          },
          () => {
            this.renderChatHistory(this.state.username, this.state.chats);
          }
        );
      }
    }
    });

    this.state.socketConn.connected(user);

    // Fetch initial activities
    this.fetchInitialActivities();

    console.log('exiting');
  };


  fetchInitialActivities = async () => {
    try {
      const response = await axios.get(`${this.state.endpoint}/activities`, { withCredentials: true });
      if (response.status === 200) {
        this.setState({ activities: response.data.activities });
      } else {
        console.error('Failed to fetch activities:', response.statusText);
      }
    } catch (error) {
      console.error('Error fetching activities:', error);
    }
  };

  // on change of input, set the value to the message state
  onChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };

  onFileChange = event => {
    this.setState({ file: event.target.files[0] });
  };

onSubmit = async e => {
  if (e.charCode === 0 && e.code === 'Enter') {
    e.preventDefault();

    let fileUrl = '';
    let message = this.state.message;

    if (this.state.file) {
      // Upload the file to S3
      const formData = new FormData();
      formData.append('file', this.state.file);
      formData.append('sender', this.state.username);
      formData.append('receiver', this.state.to);

      try {
        const response = await axios.post(
          `${this.state.endpoint}/upload`,
          formData,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
            withCredentials: true, // This should be part of the same configuration object
          }
        );

        console.log("response is: ", response);
        console.log("fileurl is : ", response.data.fileUrl);
        fileUrl = response.data.fileUrl;

        // Update the message to include the file URL
        message = `File URL is: ${fileUrl}`;
      } catch (error) {
        console.error('Error uploading file:', error);
        return;
      }
    }

    console.log("Message is : ", message);

    // Construct the message object after determining the message content
    const msg = {
      type: 'message',
      chat: {
        from: this.state.username,
        to: this.state.to,
        message: message,
        file_url: fileUrl,
        file_name: this.state.file ? this.state.file.name : '',
        file_size: this.state.file ? this.state.file.size : 0,
        file_type: this.state.file ? this.state.file.type : '',
      },
    };

    this.state.socketConn.sendMsg(msg);

    // Update local state to include the new message
    const newChat = {
      from: this.state.username,
      to: this.state.to,
      message: message,
      file_url: fileUrl,
      file_name: this.state.file ? this.state.file.name : '',
      file_size: this.state.file ? this.state.file.size : 0,
      file_type: this.state.file ? this.state.file.type : '',
    };

    this.setState(
      prevState => ({
        chats: [...prevState.chats, newChat],
        message: '',
        file: null,
        fileUrl: ''
      }),
      () => {
        this.renderChatHistory(this.state.username, this.state.chats);
      }
    );
  }
};




  getContacts = async user => {
    const res = await axios.get(
      `${this.state.endpoint}/contact-list?username=${user}`, { withCredentials: true } // get all the contacts for the given username
    );
    console.log(res.data);
    if (res.data['data'] !== undefined) {
      this.setState({ contacts: res.data.data }); // add all contacts to the contacts array
      this.renderContactList(res.data.data);
    }
  };

  fetchChatHistory = async (u1 = 'user1', u2 = 'user2') => {
    const res = await axios.get(
      `http://localhost:3000/api/chat-history?u1=${u1}&u2=${u2}`, { withCredentials: true }
    );

    console.log(res.data);
    if (res.status === 200) {
      this.setState({ chats: res.data.chats.reverse() });
      this.renderChatHistory(u1, res.data.chats.reverse());
    } else {
      console.log("sahi nahi hua");
      this.setState({ chatHistory: [] });
    }
  };

  addContact = async e => {
    e.preventDefault();
    try {
      const res = await axios.post(`${this.state.endpoint}/verify-contact`, {
        username: this.state.contact,
      });

      console.log(res.data);
      if (!res.status) {
        this.setState({ isInvalid: true });
      } else {
        // reset state on success
        this.setState({ isInvalid: false });

        let contacts = this.state.contacts;
        contacts.unshift({ // add the new contact to the contacts list
          username: this.state.contact,
          last_activity: Date.now() / 1000,
        });
        this.renderContactList(contacts);
      }
    } catch (error) {
      console.error(error);
    }
  };

  renderChatHistory = (currentUser, chats) => {
    const history = ChatHistory(currentUser, chats);
    this.setState({ chatHistory: history });
  };

  renderContactList = contacts => {
    const renderContactList = ContactList(contacts, this.sendMessageTo);

    this.setState({ renderContactList });
  };

  sendMessageTo = to => {
    this.setState({ to });
    this.fetchChatHistory(this.state.username, to);
  };

//   render() {
//     return (
//       <Container>
//         <p style={{ textAlign: 'right', paddingBottom: '10px' }}>
//           {this.state.username}
//         </p>
//         <Container paddingBottom={2}>
//           <Box>
//             <FormControl isInvalid={this.state.isInvalid}>
//               <InputGroup size="md">
//                 <Input
//                   variant="flushed"
//                   type="text"
//                   placeholder="Add Contact"
//                   name="contact"
//                   value={this.state.contact}
//                   onChange={this.onChange}
//                 />
//                 <InputRightElement width="6rem">
//                   <Button
//                     colorScheme={'teal'}
//                     h="2rem"
//                     size="lg"
//                     variant="solid"
//                     type="submit"
//                     onClick={this.addContact}
//                   >
//                     Add
//                   </Button>
//                 </InputRightElement>
//               </InputGroup>
//               {!this.state.isContactInvalid ? (
//                 ''
//               ) : (
//                 <FormErrorMessage>contact does not exist</FormErrorMessage>
//               )}
//             </FormControl>
//           </Box>
//         </Container>
//         <Flex>
//           <Box
//             textAlign={'left'}
//             overflowY={'scroll'}
//             flex="1"
//             h={'32rem'}
//             borderWidth={1}
//             borderRightWidth={0}
//             borderRadius={'xl'}
//           >
//             {this.state.renderContactList}
//           </Box>

//           <Box flex="2">
//             <Container
//               borderWidth={1}
//               borderLeftWidth={0}
//               borderBottomWidth={0}
//               borderRadius={'xl'}
//               textAlign={'right'}
//               h={'25rem'}
//               padding={2}
//               overflowY={'scroll'}
//               display="flex"
//               flexDirection={'column-reverse'}
//             >
//               {this.state.chatHistory}
//             </Container>

//             <Box flex="1">
//               <FormControl onKeyDown={this.onSubmit} onSubmit={this.onSubmit}>
//                 <Textarea
//                   type="submit"
//                   borderWidth={1}
//                   borderRadius={'xl'}
//                   minH={'7rem'}
//                   placeholder="Aur Sunao... Press enter to send..."
//                   size="lg"
//                   resize={'none'}
//                   name="message"
//                   value={this.state.message}
//                   onChange={this.onChange}
//                   isDisabled={this.state.to === ''}
//                 />
//                 <Input
//                   type="file"
//                   name="file"
//                   onChange={this.onFileChange}
//                 />
//               </FormControl>
//             </Box>
//           </Box>
//         </Flex>
//       </Container>
//     );
//   }
// }

// export default Chat;

render() {
  return (
    <Container>
      <p style={{ textAlign: 'right', paddingBottom: '10px' }}>
        {this.state.username}
      </p>
      <Container paddingBottom={2}>
        <Box>
          <FormControl isInvalid={this.state.isInvalid}>
            <InputGroup size="md">
              <Input
                variant="flushed"
                type="text"
                placeholder="Add Contact"
                name="contact"
                value={this.state.contact}
                onChange={this.onChange}
              />
              <InputRightElement width="6rem">
                <Button
                  colorScheme={'teal'}
                  h="2rem"
                  size="lg"
                  variant="solid"
                  type="submit"
                  onClick={this.addContact}
                >
                  Add
                </Button>
              </InputRightElement>
            </InputGroup>
            {!this.state.isContactInvalid ? (
              ''
            ) : (
              <FormErrorMessage>contact does not exist</FormErrorMessage>
            )}
          </FormControl>
        </Box>
      </Container>
      <Tabs>
        <TabList>
          <Tab>Chat</Tab>
          <Tab>Activity</Tab>
        </TabList>
        <TabPanels>
          <TabPanel>
            <Flex>
              <Box
                textAlign={'left'}
                overflowY={'scroll'}
                flex="1"
                h={'32rem'}
                borderWidth={1}
                borderRightWidth={0}
                borderRadius={'xl'}
              >
                {this.state.renderContactList}
              </Box>

              <Box flex="2">
                <Container
                  borderWidth={1}
                  borderLeftWidth={0}
                  borderBottomWidth={0}
                  borderRadius={'xl'}
                  textAlign={'right'}
                  h={'25rem'}
                  padding={2}
                  overflowY={'scroll'}
                  display="flex"
                  flexDirection={'column-reverse'}
                >
                  {this.state.chatHistory}
                </Container>

                <Box flex="1">
                  <FormControl onKeyDown={this.onSubmit} onSubmit={this.onSubmit}>
                    <Textarea
                      type="submit"
                      borderWidth={1}
                      borderRadius={'xl'}
                      minH={'7rem'}
                      placeholder="Aur Sunao... Press enter to send..."
                      size="lg"
                      resize={'none'}
                      name="message"
                      value={this.state.message}
                      onChange={this.onChange}
                      isDisabled={this.state.to === ''}
                    />
                    <Input type="file" name="file" onChange={this.onFileChange} />
                  </FormControl>
                </Box>
              </Box>
            </Flex>
          </TabPanel>
          <TabPanel>
            <Box p={4}>
              <VStack spacing={4} align="stretch">
                {this.state.activities.map((activity, index) => (
                  <Box key={index} p={4} borderWidth="1px" borderRadius="lg">
                    <Text>{activity.username}</Text>
                    <Text>{activity.action}</Text>
                    <Text>{new Date(activity.timestamp * 1000).toLocaleString()}</Text>
                    <Text>{activity.details}</Text>
                  </Box>
                ))}
              </VStack>
            </Box>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </Container>
  );
}
}

export default Chat;