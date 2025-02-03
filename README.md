# Parkn

Parkn is a microservice intended to be called by Twilio SMS. The service recieves a message, interprets the image using Google's VisionAPI, and stores a structure in MongoDB. From there, a goroutine will scan MongoDB and alert the user 24hrs before it's time to move their car.

*This project is incomplete*

## Example

<img src="https://github.com/user-attachments/assets/9d51c964-dcd5-490b-9165-86cea45356c0" alt="example" width="400"/>

## Usage

This app is intended for cloud deployment alongside an instance of MongoDB.

The POST endpoint should act as a webhook to a SMS Twilio Client.

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.
