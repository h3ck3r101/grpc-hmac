import grpclib
import hashlib
import hmac
import sys
import asyncio

from example_grpc import ExampleServiceStub
from example_pb2 import HelloRequest

# Your shared secret
SECRET_KEY = b'root'

def generate_hmac(message):
    return hmac.new(SECRET_KEY, message.decode().encode(), hashlib.sha256).hexdigest()

async def make_grpc_request():
    channel = grpclib.client.Channel('localhost', 50051)
    print("Grpc channel connected")
    stub = ExampleServiceStub(channel)
    
    # Construct your gRPC request
    request = HelloRequest(name='Alice')
    
    # Serialize the request
    serialized_request = request.SerializeToString()
    
    # Generate HMAC for the serialized request
    hmac_digest = generate_hmac(serialized_request)
    
    # Send the request with HMAC
    response = await stub.SayHello(request, metadata=({'hmac': hmac_digest}))
    
    if response:
        print("Response received:", response)

asyncio.run(make_grpc_request())