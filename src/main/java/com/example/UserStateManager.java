package com.example;

import java.util.concurrent.ConcurrentHashMap;

import java.util.ArrayList;
import java.util.List;

public class UserStateManager{
	private static final ConcurrentHashMap<Long, UserState> userStates = new ConcurrentHashMap<>();

	public static UserState getUserState (long userId){
		return userStates.computeIfAbsent(userId, id -> new UserState());
	}
	public static void removeUserState(long userId){
		userStates.remove(userId);
	}
}

class UserState{
	boolean isChoosing = false;
	boolean isChecking = false;
	boolean isSetting = false;
	boolean isActive = false;
	int currentTask = 0;
	int currentPage = 0;
	int answer = 0;
	String outputAnswer = null;
	String lastQuestion = null;
	String explanations = null;
	List<String> allWords = new ArrayList<>();
}
