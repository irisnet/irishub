# Future Improvements

Inter-Blockchain Communication will need to develop its own Message types that allow MTs to be transferred across chains. Making sure that spec is able to support the MTs created by this module should be easy. What might be more complicated is a transfer that includes optional tokenData so that a receiving chain has the option of parsing and storing it instead of making IBC queries when that data needs to be accessed (assuming that information stays up to date).
