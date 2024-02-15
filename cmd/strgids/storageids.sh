#
# Script will need to
# INPUT 
#   Given a List of identities (yikes a string list inout to teh scritp I expect)
#   Given a KeyVault Name  that the Service Principal can read
#
# 2 - Information Gathering
# Identify the list of storage appliances
# For Each Storage APpliance , [ only one at this point, so dont sweat this]
#   Identity the credentials by :
#   Read the Storage Appliance Resource
#   IDentify the Admin Credentials refernces password secret
#   Get the Secret 
#   Get the Cluster name from the StorageAppliance CR
#
# 3 - Gather the storage service
#     Identify the IP Address
#
# 4 - Identify teh Service PRincipal information for the Cluster
#
# 5 - Validate access to the keyvault
#     IF cannot then fail and say that the serviceprinciapl provided to teh cluster doesnt have access to the keyvault
#
# 6 - Iterate thorugh the list of user identities provided#
# For each
#    generate a password
#    Store in a hash temporarily
#
# 
# 7 - Store the password in keyvault
# Iterate thorugh the list of user identities provided
# For each create or update an entry in the provided keyvault 
#   following the name criteria of the cluster manger stored identities 
#   5ffad143-8f31-4e1e-b171-fa1738b14748-op1-cluster-op1-f7joqvkyph4xg-console-fcb855f5
#   At min we will use :
#   <cluster name>_user
# We will create/update the Secret in the keyvault.
#
# Make sure this suceeded  can be done, before we update the user perhaps  ****
# Can keyvault update we rollback to the previous version?  IF so that is the answer
#
# 8 - Create/Update Identities
# ssh onto the server
# Iterate thorugh the list of user identities provided
#   pureadduser if it doesnt exist
#   otherwise just change the passwd
#  If 8 failed, then we need to rollback the keyvault updates, however we should retry this a few times before 
#  we upgrade.
#