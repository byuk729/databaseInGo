class Parse:
    JSON_QUOTE = '"'
    JSON_CHARS = ['{', '}', '[', ']', ',', ':']
    TRUE_LEN = 4
    FALSE_LEN = 4
    NULL_LEN = 4
    JSON_LEFTBRACKET = '['
    JSON_RIGHTBRACKET = ']'
    JSON_LEFTBRACE = '{'
    JSON_RIGHTBRACE = '}'
    JSON_COMMA = ","
    JSON_COLON = ":"

    @staticmethod
    def lex_string(string):
        json_string = ''
        if string[0] == Parse.JSON_QUOTE:
            string = string[1:]
        else:
            return None, string
        for char in string:
            if char == Parse.JSON_QUOTE:
                return json_string, string[len(json_string)+1:]
            else:
                json_string += char
        raise ValueError(f"{string} is an invalid string: no closing quote")

    @staticmethod
    def lex_int(string):
        json_number = ''
        num_chars = [str(i) for i in range(0,10)] + ["e", ".", "-"]
        for char in string:
            if char in num_chars:
                json_number += char
            else:
                break
        if json_number == '':
            return None, string

        rest = string[len(json_number)::]
        if '.' in json_number:
            return float(json_number), rest
        
        return int(json_number), rest
        
    @staticmethod
    def lex_bool(string):
        if len(string) >= Parse.TRUE_LEN and string[:Parse.TRUE_LEN] == "true":
            return True, string[Parse.TRUE_LEN:]
        elif len(string) >= Parse.FALSE_LEN and string[:Parse.FALSE_LEN] == "false":
            return False, string[Parse.FALSE_LEN:]
        else:
            return None, string
    
    @staticmethod
    def lex_null(string):
        if len(string) >= Parse.NULL_LEN and string[:Parse.NULL_LEN] == "null":
            return None, string[Parse.TRUE_LEN:]
        else:
            return None, string

    @staticmethod
    def lex(string):
        tokens = []
        while len(string):
            json_string, string = Parse.lex_string(string)
            if json_string is not None:
                tokens.append(json_string)
                continue
            json_int, string = Parse.lex_int(string)
            if json_int is not None:
                tokens.append(json_int)
                continue
            json_bool, string = Parse.lex_bool(string)
            if json_bool is not None:
                tokens.append(json_bool)
                continue
            json_null, string = Parse.lex_null(string)
            if json_null is not None:
                tokens.append(json_null)
                continue
            if string[0] in Parse.JSON_CHARS:
                tokens.append(string[0])
                string = string[1:]
            elif string[0] == " ":
                string = string[1:]
            else:
                raise ValueError(f"{string} is not valid")
        return tokens
    
    @staticmethod
    def parse_array(tokens):
        array = []
        curToken = tokens[0]
        if curToken == Parse.JSON_RIGHTBRACKET:
            return array, tokens[1:]
        while True:
            json, tokens = Parse.parse(tokens)
            array.append(json)
            curToken = tokens[0]
            if curToken == Parse.JSON_RIGHTBRACKET:
                return array, tokens[1:]
            elif curToken != Parse.JSON_COMMA:
                raise ValueError("array is invalid, missing comma")
            else:
                tokens = tokens[1:]

    @staticmethod
    def parse_object(tokens):
        object = {}
        curToken = tokens[0]
        if curToken == Parse.JSON_RIGHTBRACE:
            return object, tokens[1:]
        while True:
            key = tokens[0]
            if not isinstance(key, str):
                raise Exception(f'Expected string key, got {key!r}')
            tokens = tokens[1:]
            if not tokens or tokens[0] != Parse.JSON_COLON:
                raise ValueError("Json objects needs colon following key")
            tokens = tokens[1:]
            value, tokens = Parse.parse(tokens)
            object[key] = value
            curToken = tokens[0]
            if curToken == Parse.JSON_RIGHTBRACE:
                return object, tokens[1:]
            elif curToken != Parse.JSON_COMMA:
                raise ValueError("object is invalid, missing comma")
            else:
                tokens = tokens[1:]
                print(tokens)

    @staticmethod  
    def parse(tokens):
        curToken = tokens[0]
        if curToken == Parse.JSON_LEFTBRACKET:
            return Parse.parse_array(tokens[1:])
        elif curToken == Parse.JSON_LEFTBRACE:
            return Parse.parse_object(tokens[1:])
        else:
            return curToken, tokens[1:]
        
    @staticmethod
    def from_string(string):
        tokens = Parse.lex(string)
        return Parse.parse(tokens)[0]
        
        



    