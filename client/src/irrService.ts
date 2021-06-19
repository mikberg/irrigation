import { IrrigationClient } from '@irrigation/protobuf/irrigation_grpc_web_pb';

const irrService = new IrrigationClient('https://irr.mikbe.com');

export default irrService;
