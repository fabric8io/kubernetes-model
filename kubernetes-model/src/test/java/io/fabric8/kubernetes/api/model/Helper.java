package io.fabric8.kubernetes.api.model;

import java.io.IOException;
import java.io.InputStream;
import java.util.Scanner;

public class Helper {

    public static String loadJson(String path) {
        try (InputStream resourceAsStream = Helper.class.getResourceAsStream(path)) {
            final Scanner scanner = new Scanner(resourceAsStream).useDelimiter("\\A");
            return scanner.hasNext() ? scanner.next() : "";
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
