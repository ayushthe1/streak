// this issue is fixed. thanks. I have another issue. user can upload  file and send it to other users but currently this is not working. here is the relevant section of code in Chat component. If you want any more info then ask.

import React, { Component } from 'react';
import axios from 'axios';
import SearchBar from '../Searchbar';

import SocketConnection from '../../socket-connection';
import ActivityPage from '../Activitypage';

import {
  Container,
  Flex,
  Textarea,
  Box,
  FormControl,
  FormErrorMessage,
  ModalFooter,
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
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalCloseButton,
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
      endpoint: 'https://api.ayushsharma.co.in/api',
      contact: '',
      contacts: [],
      renderContactList: [],
      chats: [],
      chatHistory: [],
      activities: [],
      msgs: [],
      file: null,
      fileUrl: '',
      searchResults: [],
      isSearchModalOpen: false,
      selectedContact: '',
      selectedFileName: '',
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
      console.log('Message is :', msg);

      if (msg.type === 'activity') {
        this.setState(prevState => ({
          activities: [msg, ...prevState.activities],
        }));
      } else {
        // update UI only when message is between from and to
        if (
          this.state.username === msg.to ||
          this.state.username === msg.from
        ) {
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

  handleSearchResults = response => {
    const results = response.result.data.Get.Users;
    this.setState({ searchResults: results, isSearchModalOpen: true });
  };

  closeSearchModal = () => {
    this.setState({ isSearchModalOpen: false });
  };

  fetchInitialActivities = async () => {
    try {
      const response = await axios.get(`${this.state.endpoint}/activities`, {
        withCredentials: true,
      });
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
    const file = event.target.files[0];
    console.log('File selected:', file);
    this.setState({
      file: file,
      selectedFileName: file ? file.name : '',
    });
  };

  onSubmit = async e => {
    if (
      (e.type === 'keydown' && e.charCode === 0 && e.code === 'Enter') ||
      e.type === 'click'
    ) {
      e.preventDefault();

      let fileUrl = '';
      let fileName = '';
      let fileSize = 0;
      let fileType = '';

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

          console.log('response is: ', response);
          console.log('fileurl is : ', response.data.fileUrl);
          fileUrl = response.data.fileUrl;
          fileName = this.state.file.name;
          fileSize = this.state.file.size;
          fileType = this.state.file.type;

          // Update the message to include the file URL
          // message = `File URL is: ${fileUrl}`;
        } catch (error) {
          console.error('Error uploading file:', error);
          return;
        }
      }

      // console.log('Message is : ', message);

      // Construct the message object after determining the message content
      const msg = {
        type: 'message',
        chat: {
          from: this.state.username,
          to: this.state.to,
          message: this.state.message,
          file_url: fileUrl,
          file_name: fileName,
          file_size: fileSize,
          file_type: fileType,
        },
      };

      this.state.socketConn.sendMsg(msg);

      // Update local state to include the new message
      const newChat = {
        from: this.state.username,
        to: this.state.to,
        // message: message,
        message: this.state.message, // Keep the original message
        file_url: fileUrl,
        file_name: fileName,
        file_size: fileSize,
        file_type: fileType,
        timestamp: Math.floor(Date.now() / 1000),
      };
      document.getElementById('file-upload').value = '';
      this.setState(
        prevState => ({
          chats: [...prevState.chats, newChat],
          message: '',
          file: null,
          fileUrl: '',
          selectedFileName: '',
        }),
        () => {
          this.renderChatHistory(this.state.username, this.state.chats);
        }
      );
    }
  };

  getContacts = async user => {
    const res = await axios.get(
      `${this.state.endpoint}/contact-list?username=${user}`,
      { withCredentials: true }
    );
    console.log(res.data);
    if (res.data['data'] !== undefined) {
      this.setState({ contacts: res.data.data }, () => {
        this.renderContactList(this.state.contacts);
      });
    }
  };

  fetchChatHistory = async (u1 = 'user1', u2 = 'user2') => {
    const res = await axios.get(
      `https://api.ayushsharma.co.in/api/chat-history?u1=${u1}&u2=${u2}`, { withCredentials: true }
    );

    console.log(res.data);
    if (res.status === 200) {
      this.setState({ chats: res.data.chats.reverse() });
      this.renderChatHistory(u1, res.data.chats.reverse());
    } else {
      console.log('sahi nahi hua');
      this.setState({ chatHistory: [] });
    }
  };

  renderChatHistory = (currentUser, chats) => {
    const history = ChatHistory(currentUser, chats);
    this.setState({ chatHistory: history });
  };

  renderContactList = contacts => {
    const renderContactList = ContactList(
      contacts,
      this.sendMessageTo,
      this.state.selectedContact
    );
    this.setState({ renderContactList });
  };

  sendMessageTo = to => {
    this.setState({ to, selectedContact: to }, () => {
      this.renderContactList(this.state.contacts);
      this.fetchChatHistory(this.state.username, to);
    });
  };

  render() {
    const { searchResults, isSearchModalOpen, activities } = this.state;
    return (
      <Box bg="gray.900" minHeight="100vh" p={5}>
        <Flex direction="column" h="100vh" maxW="1200px" mx="auto">
          <Flex justify="space-between" align="center" mb={4}>
            <SearchBar
              from={this.state.username}
              onSearchResults={this.handleSearchResults}
            />
            <Text color="purple.300" fontWeight="bold">
              {this.state.username}
            </Text>
          </Flex>
  
          <Tabs isFitted variant="enclosed" colorScheme="purple" flex={1}>
            <TabList>
              <Tab>Chat</Tab>
              <Tab>Activities</Tab>
            </TabList>
            <TabPanels>
              <TabPanel p={0}>
                <Flex
                  flex={1}
                  borderRadius="xl"
                  overflow="hidden"
                  boxShadow="xl"
                  h="calc(100vh - 150px)"
                >
                  <Box
                    w={{ base: '100%', md: '300px' }}
                    bg="gray.700"
                    overflowY="auto"
                    borderRightWidth={1}
                    borderColor="gray.600"
                    display={{
                      base: this.state.to ? 'none' : 'block',
                      md: 'block',
                    }}
                  >
                    <Box p={4}>
                      {/* <FormControl isInvalid={this.state.isInvalid} mb={4}> */}
                        {/* <InputGroup size="md">
                          <Input
                            bg="gray.600"
                            color="white"
                            border="none"
                            placeholder="Add Contact"
                            name="contact"
                            value={this.state.contact}
                            onChange={this.onChange}
                          /> */}
                          {/* <InputRightElement width="4.5rem">
                            <Button
                              h="1.75rem"
                              size="sm"
                              colorScheme="purple"
                              onClick={this.addContact}
                            >
                              Add
                            </Button>
                          </InputRightElement> */}
                        {/* </InputGroup> */}
                        {/* {this.state.isContactInvalid && (
                          <FormErrorMessage>
                            Contact does not exist
                          </FormErrorMessage>
                        )} */}
                      {/* </FormControl> */}
                      {this.state.renderContactList}
                    </Box>
                  </Box>
  
                  <Flex
                    direction="column"
                    flex={1}
                    bg="gray.800"
                    display={{
                      base: this.state.to ? 'flex' : 'none',
                      md: 'flex',
                    }}
                  >
                    <Box flex={1} overflowY="auto" p={4}>
                      {this.state.chatHistory}
                    </Box>
  
                    <Box p={4} bg="gray.700">
                      {this.state.to !== '' ? (
                        <FormControl
                          onKeyDown={this.onSubmit}
                          onSubmit={this.onSubmit}
                        >
                          <Textarea
                            bg="gray.600"
                            color="white"
                            border="none"
                            borderRadius="md"
                            placeholder="Type your message here... Press enter to send"
                            _placeholder={{ color: 'gray.400' }}
                            mb={2}
                            name="message"
                            value={this.state.message}
                            onChange={this.onChange}
                            rows={3}
                          />
                          <Flex
                            justify="space-between"
                            align="center"
                            flexWrap="wrap"
                          >
                            <Input
                              type="file"
                              name="file"
                              onChange={this.onFileChange}
                              hidden
                              id="file-upload"
                            />
                            <Button
                              as="label"
                              htmlFor="file-upload"
                              colorScheme="purple"
                              size="sm"
                              mb={{ base: 2, sm: 0 }}
                            >
                              Attach File
                            </Button>
                            {this.state.selectedFileName && (
                              <Text
                                fontSize="sm"
                                color="gray.400"
                                ml={2}
                                mb={{ base: 2, sm: 0 }}
                              >
                                {this.state.selectedFileName}
                              </Text>
                            )}
                            <Button
                              colorScheme="purple"
                              size="sm"
                              onClick={this.onSubmit}
                            >
                              Send
                            </Button>
                          </Flex>
                        </FormControl>
                      ) : (
                        <Flex justify="center" align="center" h="100%">
                          <Text color="gray.400">
                            Select a contact to start chatting
                          </Text>
                        </Flex>
                      )}
                    </Box>
                  </Flex>
                </Flex>
              </TabPanel>
              <TabPanel>
                <ActivityPage activities={activities} />
              </TabPanel>
            </TabPanels>
          </Tabs>
  
          <Modal isOpen={isSearchModalOpen} onClose={this.closeSearchModal}>
            <ModalOverlay />
            <ModalContent>
              <ModalHeader>Search Results</ModalHeader>
              <ModalCloseButton />
              <ModalBody>
                {searchResults.length === 0 ? (
                  <Text>No results found.</Text>
                ) : (
                  <VStack spacing={4} align="stretch">
                    {searchResults.map((result, index) => (
                      <Box
                        key={index}
                        p={3}
                        shadow="md"
                        borderWidth="1px"
                        borderRadius="md"
                      >
                        <Flex direction="column">
                          <Flex justifyContent="space-between">
                            <Text fontWeight="bold">From: {result.from}</Text>
                            <Text fontWeight="bold">To: {result.to}</Text>
                          </Flex>
                          <Text mt={2}>{result.message}</Text>
                        </Flex>
                      </Box>
                    ))}
                  </VStack>
                )}
              </ModalBody>
              <ModalFooter>
                <Button
                  colorScheme="blue"
                  mr={3}
                  onClick={this.closeSearchModal}
                >
                  Close
                </Button>
              </ModalFooter>
            </ModalContent>
          </Modal>
        </Flex>
      </Box>
    );
  }
}

export default Chat;
