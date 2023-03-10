from hfc.fabric import Client

import asyncio

loop = asyncio.get_event_loop()

cli = Client(net_profile="./network_2_0.json")

# print(cli.organizations)  # orgs in the network
# print(cli.peers)  # peers in the network
# print(cli.orderers)  # orderers in the network
# print(cli.CAs)  # ca nodes in the network

org1_admin = cli.get_user(org_name='org1.example.com', name='Admin') # User instance with the Org1 admin's certs
org2_admin = cli.get_user(org_name='org2.example.com', name='Admin') # User instance with the Org2 admin's certs
orderer_admin = cli.get_user(org_name='orderer.example.com', name='Admin') # User instance with the orderer's certs

response = loop.run_until_complete(cli.channel_create(
            orderer='orderer.example.com',
            channel_name='businesschannel',
            requestor=org1_admin,
            config_yaml='./configtx/',
            ))

# response = true is returned if the channel is created successfully
print(response)