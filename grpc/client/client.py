import sys
import grpc
import goraft_pb2
import goraft_pb2_grpc


def main():
    if len(sys.argv) < 3:
        print("Usage: python client.py <host:port> <command>")
        print("  Commands: SET:key:value  GET:key  DEL:key")
        sys.exit(1)

    target = sys.argv[1]
    command = sys.argv[2]

    try:
        channel = grpc.insecure_channel(target)
        stub = goraft_pb2_grpc.NodeStub(channel)

        request = goraft_pb2.CommandRequest(command=command)
        response = stub.ExecuteCommand(request, timeout=3)

        if response.success:
            print(f"OK: {response.message}")
        else:
            msg = response.message
            print(f"FAIL: {msg}")

    except grpc.RpcError as err:
        code = err.code()
        details = err.details()

        if code == grpc.StatusCode.UNAVAILABLE:
            print(f"FAIL: node unavailable at {target}")
            print("Reason: could not connect to the node. Is it running?")
        elif code == grpc.StatusCode.DEADLINE_EXCEEDED:
            print(f"FAIL: request to {target} timed out")
        else:
            print(f"FAIL: RPC error while contacting {target}")
            print(f"Code: {code.name}")
            print(f"Details: {details}")

        sys.exit(1)


if __name__ == "__main__":
    main()