import java.util.concurrent.ConcurrentHashMap;

public class UserStateManager{
	private static final ConcurrentHashMap<Long, UserState> userStates = new ConcurrentHashMap<>();

	public static UserState getUserState (long userId){
		return userStates.computeIfAbsent(userId, id -> new UserState());
	}
}

class UserState{
	boolean isChoosing = true;
	boolean isChecking = false;
	int currentTask = 0;
	String lastQuestion = null;
	String explanations = null;
}
