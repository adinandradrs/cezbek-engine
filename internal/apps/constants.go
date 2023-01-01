package apps

const DefaultTrxId = "11111"
const StatusActive = 1
const StatusInactive = 0
const SuccessCode = "8000"
const SuccessMsgOk = "OK"
const SuccessMsgSubmit = "Data submitted successfully"
const SuccessMsgDataFound = "Here is your data"

const ErrCodeSomethingWrong = "9001"
const ErrMsgSomethingWrong = "Something went wrong, please contact administrator"
const ErrMsgESBUnavailable = "The core engines are out of service"
const ErrCodeESBUnavailable = "9003"
const ErrMsgUnauthorized = "Unauthorized access, please try again"
const ErrCodeUnauthorized = "9005"
const ErrMsgSubmitted = "Failed to submit the data"
const ErrCodeSubmitted = "9006"
const ErrMsgTokenExpired = "The session has been expired"
const ErrCodeTokenExpired = "9007"
const ErrMsgBadPayload = "Please check again the given input"
const ErrCodeBadPayload = "9008"
const ErrMsgInvalidChannel = "The given channel is invalid, try another channel"
const ErrCodeInvalidChannel = "9009"

const ErrCodeBussPartnerExists = "BR-01"
const ErrMsgBussPartnerExists = "The given partner data is exists on system"

const HeaderClientTrxId = "x-client-trxid"
const HeaderClientChannel = "x-client-channel"
const HeaderClientDeviceId = "x-client-device"
const HeaderClientVersion = "x-client-version"
const HeaderClientRefToken = "x-client-refresh-token"
const HeaderClientOs = "x-client-os"
const HeaderClientTimestamp = "x-client-timestamp"
const HeaderClientSignature = "x-client-signature"
const HeaderApiKey = "x-api-key"
const ChannelB2BClient = "B2BCLIENT"
const ChannelEBizKezbek = "EBIZKEZBEK"
