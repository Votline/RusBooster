import java.io.File;
import java.io.IOException;
import java.util.Map;
import java.util.HashMap;

import com.fasterxml.jackson.databind.ObjectMapper;

public class Settings{
	private Map<String, String> exercise = new HashMap<>();
	private ObjectMapper objMapper = new ObjectMapper();
	private File base = new File("settings.json");
	private String message;

	public void chooseExercise(String userId, String numberOfExercise){
		try{
			if(base.exists() && base.length() > 0 ){
				exercise = objMapper.readValue(base, Map.class);

				exercise.put(userId, numberOfExercise);
				objMapper.writeValue(base, exercise);

				message = "Добавил: " + userId + " " + numberOfExercise;
				System.out.println(message);
			}
		}
		catch(IOException e){
			e.printStackTrace();
		}
	}
}
