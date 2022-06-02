# goCheck

This is a tool I built for a task at my job. It reads an .xml file with the following structure:

```xml
<rules>
    <rule name="..." stopProcessing="...">
        <match url="regex" />
        <conditions logicalGrouping="..." trackAllCaptures="...">
            <add input="..." pattern="link that will be redirected" />
            <add input="..." pattern="link that will be redirected" />
        </conditions>
        <action type="..." url="where pattern will be redirected to" redirectType="..."/>
    </rules>
    ...
</rules>
```

Where `pattern` is a link that will be accessed by an user and `url` is where he needs to be redirected. The program will test if all the redirects are working correctly by making a request to them and then store the results in a file called `results.csv`.

## Usage

Use `timeout` to specify how much time the client should wait for a response.
Use `filename` to specify the .xml file you wish to parse and check.
Use `filter` to filter the output file to show only files that received a response which was not expected (when the url found is different from the expected url - `url`)

