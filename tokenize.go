package dtsh

// There are two types of tokens that dtsh needs to be aware of: regular tokens,
// which are either words or strings delimited by double quotes, and literal
// tokens, which are strings delimited by single quotes. The reason we need to
// distinguish between them is that variable substitution is not allowed inside
// literal tokens.
type TokenType int

const (
    TokenRegular TokenType = iota
    TokenLiteral
)

// Tokens are just a string and a TokenType.
type Token struct {
    Type  TokenType
    Value string
}

// Write out the string inside the Token.
func (t Token) String() string {
    return t.Value
}

// Tokenize takes an input and splits it into tokens. It's modeled as a state
// machine with five states.
//
//  stateWhitespace -> (double quote) -> stateString
//                     (single quote) -> stateLiteral
//                     (whitespace)   -> stateWhitespace
//                     (character)    -> character pushed
//                                       stateWord
//
//  stateWord -> (whitespace)   -> token pushed and cleared
//                                 stateWhitespace
//               (double quote) -> stateString
//               (single quote) -> stateLiteral
//               (character)    -> character pushed
//
//  stateString -> (double quote) -> token pushed and cleared
//                                   stateWhitespace (?)
//                 (backslash)    -> save state to lastState
//                                   stateBackslash
//                 (character)    -> character pushed
//
//  stateLiteral -> (single quote) -> token pushed and cleared
//                                    stateWhitespace (?)
//                  (backslash)    -> save state to lastState
//                                    stateBackslash
//                  (character)    -> character pushed
//
//  stateBackslash -> n           -> \n pushed
//                                   lastState
//                    r           -> \r pushed
//                                   lastState
//                    t           -> \t pushed
//                                   lastState
//                    b           -> \b pushed
//                                   lastState
//                    f           -> \f pushed
//                                   lastState
//                    v           -> \v pushed
//                                   lastState
//                    (character) -> character pushed
//                                   lastState
//
// TODO: UTF-8 code points (\x..)
func Tokenize(s string) []Token {
    type stateT int
    const (
        stateWhitespace stateT = iota
        stateWord
        stateString
        stateLiteral
        stateBackslash
    )

    // The token we're working on is stored as a rune slice. We convert it
    // to a Token only when pushing it to the tokens slice.
    var tokens []Token
    var token []rune
    var state stateT = stateWhitespace
    var lastState stateT

    // Loop over the given string's characters.
    // Add a space after the string, to prevent the state machine from not
    // pushing the last token.
    for _, char := range s + " " {
        switch state {
        case stateWhitespace:
            switch char {
            case '"':
                state = stateString
            case '\'':
                state = stateLiteral
            case ' ':
                state = stateWhitespace
            default:
                token = append(token, char)
                state = stateWord
            }
        case stateWord:
            switch char {
            case ' ':
                tokens = append(tokens, Token{
                    Type:  TokenRegular,
                    Value: string(token),
                })
                token = nil
                state = stateWhitespace
            case '"':
                state = stateString
            case '\'':
                state = stateLiteral
            default:
                token = append(token, char)
            }
        case stateString:
            switch char {
            case '"':
                tokens = append(tokens, Token{
                    Type:  TokenRegular,
                    Value: string(token),
                })
                token = nil
                state = stateWhitespace
            case '\\':
                lastState = state
                state = stateBackslash
            default:
                token = append(token, char)
            }
        case stateLiteral:
            switch char {
            case '\'':
                tokens = append(tokens, Token{
                    Type:  TokenLiteral,
                    Value: string(token),
                })
                token = nil
                state = stateWhitespace
            case '\\':
                lastState = state
                state = stateBackslash
            default:
                token = append(token, char)
            }
        case stateBackslash:
            switch char {
            case 'n':
                token = append(token, '\n')
            case 'r':
                token = append(token, '\r')
            case 't':
                token = append(token, '\t')
            case 'b':
                token = append(token, '\b')
            case 'f':
                token = append(token, '\f')
            case 'v':
                token = append(token, '\v')
            default:
                token = append(token, char)
            }
            state = lastState
        }
    }
    return tokens
}
