import unittest
from scripts.parsing import Parse

lex_string = Parse.lex_string
lex = Parse.lex
class TestParsing(unittest.TestCase):
    def test_normal_string_parses_correctly(self):
        # Assertions check if expected results match actual results
        json_string, restOfString = lex_string('"foo"6767676767')
        self.assertEqual(json_string, "foo")
        self.assertEqual(restOfString, "6767676767")
    def test_broken_string_raises_error(self):
        # Assertions check if expected results match actual results
        with self.assertRaises(ValueError) as context:
            lex_string('"foo6767676767')
    def test_normal_json_parses_correctly(self):
        self.assertEqual(lex('{"foo": [1, 2, {"bar": 2}]}'),
                ['{', 'foo', ':', '[', 1, ',', 2, ',', '{', 'bar', ':', 2, '}', ']', '}'])

    def test_array_without_objects_parses_to_python_list(self):
        self.assertEqual(
            Parse.from_string('[1, 2, [3, 4], "five", true]'),
            [1, 2, [3, 4], 'five', True]
        )
    
    def test_object_without_arrays_parses_to_python_dict(self):
        self.assertEqual(
            Parse.from_string('{"name": "Bryan", "age": 21, "active": true, "address": {"city": "Charlottesville"}}'),
            {'name': 'Bryan', 'age': 21, 'active': True, 'address': {'city': 'Charlottesville'}}
        )

    def test_normal_json_parses_to_python_object(self):
        self.assertEqual(
            Parse.from_string('{"foo": [1, 2, {"bar": 2}]}'),
            {'foo': [1, 2, {'bar': 2}]}
        )

if __name__ == '__main__':
    unittest.main()
    maxDiff = None