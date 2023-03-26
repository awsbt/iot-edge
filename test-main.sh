#!/usr/bin/env bash
clear
amURL=https://fr0.zte.peraton.com/am
thingId=pem-bot-01
#thingId=thingymabot
tree=RegisterThings
#keyfile=$(pwd)/examples/resources/eckey1.key.pem
#certfile=$(pwd)/examples/resources/thingymabot.cert.pem
keyfile=$(pwd)/examples/resources/pem-bot-01.key.pem
certfile=$(pwd)/examples/resources/pem-bot-01.crt
#keyfile=$(pwd)/examples/resources/DINPE-Bot-015-PrivateKey.pem
#certfile=$(pwd)/examples/resources/DINPE-Bot-015.crt

# Initiate the authentication request:
authCallback=$(curl \
    --silent \
    --header 'Accept-API-Version: resource=2.0, protocol=1.0' \
    --header 'Content-Type: application/json' \
    --request POST \
    "$amURL/json/authenticate?authIndexType=service&authIndexValue=$tree")

#echo "authCallback="$authCallback
#read -p "Press any key to resume ..."

# Extract challenge:
challenge=$(echo "$authCallback" | \
    jq ".callbacks[0].output[0].value")

#echo "challenge="$challenge
#read -p "Press any key to resume ..."

# Create the signed JWT for the Authenticate Thing Node:
signedJWT=$(go run ./cmd/auth-jwt -a "/" -s "$thingId" -c "$challenge" --key "$keyfile")

#echo "signedJWT="$signedJWT
#read -p "Press any key to resume ..."

# Modify callback:
authCallback=$(echo "$authCallback" | \
    jq ".callbacks[0].input[0].value = \"$signedJWT\"")
	
#echo "modify authCallback="$authCallback
#read -p "Press any key to resume ..."

# Complete the authentication request:
authResponse=$(curl \
    --silent \
    --header 'Accept-API-Version: resource=2.0, protocol=1.0' \
    --header 'Content-Type: application/json' \
    --request POST \
    --data "$authCallback" \
    "$amURL/json/authenticate?authIndexType=service&authIndexValue=$tree")
	
#echo "authResponse="$authResponse
#read -p "Press any key to resume ..."

ssoToken=$(jq -r '.tokenId' <(echo $authResponse))

#echo "ssoToken="$ssoToken
#read -p "Press any key to resume ..."

if [ "$ssoToken" != "null" ]; then
    #echo "Authentication complete $ssoToken"

    jwt=$(go run ./cmd/things-jwt \
    -u "$amURL/json/things/*?_action=get_access_token&realm=/" \
    -k "$keyfile" \
    --custom '{"scope":["publish"]}')

    #echo "jwt="$jwt

    accessTokenResponse=$(curl \
    --silent \
    --header 'accept-api-version: protocol=2.0,resource=1.0' \
    --header 'content-type: application/jose' \
    --cookie "iPlanetDirectoryPro=$ssoToken" \
    --request POST \
    --data "$jwt" \
    "$amURL/json/things/*?_action=get_access_token&realm=/")

    accessToken=$(echo "$accessTokenResponse" | jq -r '.access_token')

    #echo "accessTokenResponse="$accessTokenResponse

    echo "accessToken="$accessToken

    exit 0

else
    callbackId=$(echo "$authResponse" | \
    jq '[ .callbacks[0].output[] | select( .name | contains("id")) ]' | \
     jq  .[0].value)
    if [ "$callbackId" = '"jwt-pop-registration"' ]; then
        echo "Thing is unknown to AM, please continue to Registration"
    else
        echo "Something has gone wrong"
        exit 1
    fi
fi

#read -p "Press any key to resume ..."

# Extract challenge:
challenge=$(echo "$authResponse" | \
    jq ".callbacks[0].output[0].value")

echo "challenge="$challenge
#read -p "Press any key to resume ..."

# Create the signed registration JWT for the Register Thing Node:
signedJWT=$(go run ./cmd/auth-jwt -a "/" -s "$thingId" -c "$challenge" --key "$keyfile" --certificate $certfile)

echo "signedJWT="$signedJWT
#read -p "Press any key to resume ..."

# Modify callback:
regCallback=$(echo "$authResponse" | \
        jq ".callbacks[0].input[0].value = \"$signedJWT\"")

#echo "modify regCallback="$regCallback
#read -p "Press any key to resume ..."

# Complete the registration request:
regResponse=$(curl \
    --silent \
    --header 'Accept-API-Version: resource=2.0, protocol=1.0' \
    --header 'Content-Type: application/json' \
    --request POST \
    --data "$regCallback" \
    "$amURL/json/authenticate?authIndexType=service&authIndexValue=$tree")

echo "regResponse="$regResponse
#read -p "Press any key to resume ..."

ssoToken=$(jq -r '.tokenId' <(echo $regResponse))
#echo "ssoToken="$ssoToken
#echo "{ssoToken}=${ssoToken}"

jwt=$(go run ./cmd/things-jwt \
-u "$amURL/json/things/*?_action=get_access_token&realm=/" \
-k "$keyfile" \
--custom '{"scope":["publish"]}')

echo "jwt="$jwt

accessTokenResponse=$(curl \
--silent \
--header 'accept-api-version: protocol=2.0,resource=1.0' \
--header 'content-type: application/jose' \
--cookie "iPlanetDirectoryPro=$ssoToken" \
--request POST \
--data "$jwt" \
"$amURL/json/things/*?_action=get_access_token&realm=/")

accessToken=$(echo "$accessTokenResponse" | jq -r '.access_token')

echo "accessTokenResponse="$accessTokenResponse

echo "accessToken="$accessToken


