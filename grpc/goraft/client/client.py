import sys
import grpc
import goraft_pb2
import goraft_pb2_grpc


def main():
    if len(sys.argv) < 3:
        print(f"Usage: python client.py <host:port> <command>")
        print(f"  Commands: SET:key:value  GET:key  DEL:key")
        sys.exit(1)

    target = sys.argv[1]
    command = sys.argv[2]

    channel = grpc.insecure_channel(target)
    stub = goraft_pb2_grpc.NodeStub(channel)

    request = goraft_pb2.CommandRequest(command=command)
    response = stub.ExecuteCommand(request)

    if response.success:
        print(f"OK: {response.message}")
    else:
        msg = response.message
        if response.leader_address:
            msg += f" (leader: {response.leader_address})"
        print(f"FAIL: {msg}")


if __name__ == "__main__":
    main()
