import React from 'react';

import { Box, Heading, Center } from '@chakra-ui/react';

function Footer() {
  return (
    <Box padding={8}>
      <Center>
        <Heading size="sm">Powered by Golang</Heading>
      </Center>
      <Center>
        <Heading fontStyle={'italic'} size="sm" paddingTop={2}>
          <a href="github.com/ayushthe1" rel="noreferrer" target={'_blank'}>
            Ayush Sharma
          </a>
        </Heading>
      </Center>
    </Box>
  );
}

export default Footer;
