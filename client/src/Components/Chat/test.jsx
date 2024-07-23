// class Chat extends Component {
//     constructor(props) {
//       super(props);
//       this.state = {
//         socketConn: '',
//         username: '',
//         message: '',
//         to: '',
//         isInvalid: false,
//         endpoint: 'http://localhost:3000/api',
//         contact: '',
//         contacts: [],
//         renderContactList: [],
//         chats: [],
//         chatHistory: [],
//         activities: [],
//         msgs: [],
//         file: null,
//         fileUrl: '',
//         localStream: null,
//         remoteStream: null,
//         peerConnection: null,
//         signalingServer: null,
//       };
  
//       this.localVideoRef = createRef();
//       this.remoteVideoRef = createRef();
//     }
  
//     componentDidMount = async () => {
//       const queryParams = new URLSearchParams(window.location.search);
//       const user = queryParams.get('u');
//       this.setState({ username: user });
//       this.getContacts(user); // get all contacts of the user
  
//       const conn = new SocketConnection();
//       await this.setState({ socketConn: conn });
//       this.state.socketConn.connect(message => {
//         const msg = JSON.parse(message.data);
//         console.log("Message is :", msg);
  
//         if (msg.type === 'activity') {
//           this.setState((prevState) => ({
//             activities: [msg, ...prevState.activities],
//           }));
//         } else if (msg.type === 'offer' || msg.type === 'answer' || msg.type === 'ice-candidate') {
//           console.log("SignallingMessage received of type : ", msg.type)
//           this.handleSignalingMessage(msg);
//         } else {
//           if (this.state.username === msg.to || this.state.username === msg.from) {
//             this.setState(
//               {
//                 chats: [...this.state.chats, msg],
//               },
//               () => {
//                 this.renderChatHistory(this.state.username, this.state.chats);
//               }
//             );
//           }
//         }
//       });
  
//       this.state.socketConn.connected(user);
//       this.fetchInitialActivities();
//     };
    
//     startCall = async (targetUsername) => {
//       try {
//         const localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
//         this.localVideoRef.current.srcObject = localStream;
//         this.setState({ localStream });
  
//         const peerConnection = new RTCPeerConnection({
//           iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
//         });
  
//         peerConnection.onicecandidate = ({ candidate }) => {
//           if (candidate) {
//             this.state.socketConn.sendMsg({type:'signalling-message', signallingMessage:{
//               type: 'ice-candidate',
//               candidate,
//               targetUsername,
//             }});
//           }
//         };
  
//         peerConnection.ontrack = (event) => {
//           this.remoteVideoRef.current.srcObject = event.streams[0];
//           this.setState({ remoteStream: event.streams[0] });
//         };
  
//         localStream.getTracks().forEach((track) => peerConnection.addTrack(track, localStream));
  
//         const offer = await peerConnection.createOffer();
//         await peerConnection.setLocalDescription(offer);
  
//         this.setState({ peerConnection });
  
//         this.state.socketConn.sendMsg({type : 'signalling-message',signallingMessage:{
//           type: 'offer',
//           offer,
//           targetUsername,
//         }});
//       } catch (error) {
//         console.error('Error starting call:', error);
//       }
//     };
  
//     handleSignalingMessage = async (msg) => {
//       console.log("Message received in handleSignallingMessage is :", msg)
//       const { type, offer, answer, candidate, senderId } = msg;
//       const { peerConnection } = this.state;
  
//       switch (type) {
//         case 'offer':
//           console.log("Case offer received")
//           await this.handleOffer(offer, senderId);
//           break;
//         case 'answer':
//           console.log("Case answer received")
//           await peerConnection.setRemoteDescription(new RTCSessionDescription(answer));
//           break;
//         case 'ice-candidate':
//           console.log("Case ice-candidate received")
//           if (peerConnection != null ) {
//             await peerConnection.addIceCandidate(new RTCIceCandidate(candidate));
//           }else{
//             console.log("PEERCONNECTION IS NULL")
//           }
          
//           break;
//         default:
//           break;
//       }
//     };
  
//     handleOffer = async (offer, senderId) => {
//       const localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
//       this.localVideoRef.current.srcObject = localStream;
//       this.setState({ localStream });
  
//       const peerConnection = new RTCPeerConnection({
//         iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
//       });
  
//       peerConnection.onicecandidate = ({ candidate }) => {
//         if (candidate) {
//           this.state.socketConn.sendMsg({type : 'signalling-message',signallingMessage:{
//             type: 'ice-candidate',
//             candidate,
//             targetUsername: this.state.to,
//           }});
//         }
//       };
  
//       peerConnection.ontrack = (event) => {
//         this.remoteVideoRef.current.srcObject = event.streams[0];
//         this.setState({ remoteStream: event.streams[0] });
//       };
  
//       localStream.getTracks().forEach((track) => peerConnection.addTrack(track, localStream));
  
//       await peerConnection.setRemoteDescription(new RTCSessionDescription(offer));
//       const answer = await peerConnection.createAnswer();
//       await peerConnection.setLocalDescription(answer);
  
//       this.setState({ peerConnection });
  
//       this.state.socketConn.sendMsg({type : 'signalling-message',signallingMessage:{
//         type: 'answer',
//         answer,
//         targetUsername: this.state.to,
//       }});
//     };