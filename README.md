# twitch-acc-gen
 An account creator/follow bot for [twitch](https://twitch.tv)
 
# Reason for release
Recently I've been seeing a lot of services pop up for twitch, these services charge absolutely insane amounts and no one really knows how easy it really is to generate accounts and mass follow people. I'm hoping this program doesn't fuel the problem but instead tunes it down. I've also just been having this collecting dust on my desktop for a few months and needed a purpose for it (github stars).

# How to use
I will not be providing examples, just a general over view.

1.) get your twitch userID, to do this just inspect network traffic when you follow your account from another account. Look for "gql.twitch.tv/gql" op name "FollowButton_FollowUser".\
2.) you can use an proxy for the proxy field, I suggest [oxylabs](https://oxylabs.io)\
3.) captchaKey, currently this tool only supports [bestcaptchasolver](https://bestcaptchasolver.com), this is because, from my testing, they have the fastest response time to funcaptcha.

# Notice
Follow botting/account generation is against the TOS of twitch. I do not condone anyone using this tool. I made this tool as a passion project and decided to release it. Please do not **skid** this and please do not use this to create a toxic place on twitch.

# Credit(s)
[azaelgg](https://github.com/azaelgg) for even getting me interested in this in the first place. He also provided code to help me understand websockets better.
