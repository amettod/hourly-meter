# :book: Introduction 
The application performs syntax analysis of *.html file received by "Universal configurator of electric meters Mercury" program. On its basis, it creates *.xml files with hourly readings for their further transfer to the energy company of Moscow and the region.

# :wrench: How to use

* **Command-line interface**
    ```shellscript
    $ go build ./cmd/cli
    ```
    ```shellscript
    $ ./cli -filename="filename.html" \
        -contract="contract humber" \
        -name="company name" \
        -coefficient="power factory"
    ```
* **Web interface**
    ```shellscript
    $ go build ./cmd/web
    ```
    ```shellscript
    $ ./web 
    ```

* **Use Makefile**
    ```shellscript
    $ make
    ```